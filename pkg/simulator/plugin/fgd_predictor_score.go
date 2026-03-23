package plugin

import (
	"context"
	"fmt"
	"math"
	"sort"

	"strconv"

	// log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	simontype "github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/type"
	gpushareutils "github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/type/open-gpu-share/utils"
	"github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/utils"
)

type FGDPredictorScorePlugin struct {
	handle      framework.Handle
	typicalPods *simontype.TargetPodList
}

var _ framework.ScorePlugin = &FGDPredictorScorePlugin{}

func NewFGDPredictorScorePlugin(
	_ runtime.Object,
	handle framework.Handle,
	typicalPods *simontype.TargetPodList,
) (framework.Plugin, error) {
	plugin := &FGDPredictorScorePlugin{
		handle:      handle,
		typicalPods: typicalPods,
	}
	allocateGpuIdFunc[plugin.Name()] = allocateGpuIdBasedOnFGDPredictorScore
	return plugin, nil
}

func (p *FGDPredictorScorePlugin) Name() string {
	return simontype.FGDPredictorScorePluginName
}

func (p *FGDPredictorScorePlugin) Score(
	ctx context.Context,
	state *framework.CycleState,
	pod *v1.Pod,
	nodeName string,
) (int64, *framework.Status) {
	if podReq, _ := resourcehelper.PodRequestsAndLimits(pod); len(podReq) == 0 {
		return framework.MaxNodeScore, framework.NewStatus(framework.Success)
	}

	nodeResPtr := utils.GetNodeResourceViaHandleAndName(p.handle, nodeName)
	if nodeResPtr == nil {
		return framework.MinNodeScore, framework.NewStatus(
			framework.Error,
			fmt.Sprintf("failed to get nodeRes(%s)", nodeName),
		)
	}
	nodeRes := *nodeResPtr

	podRes := utils.GetPodResource(pod)
	if !utils.IsNodeAccessibleToPod(nodeRes, podRes) {
		return framework.MinNodeScore, framework.NewStatus(
			framework.Error,
			fmt.Sprintf("Node (%s) %s does not match GPU type request of pod %s",
				nodeName, nodeRes.Repr(), podRes.Repr()),
		)
	}

	score, _ := calculateOptimizedScore(nodeRes, podRes, p.typicalPods)
	return score, framework.NewStatus(framework.Success)
}

func (p *FGDPredictorScorePlugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

func calculateOptimizedScore(
	nodeRes simontype.NodeResource,
	podRes simontype.PodResource,
	typicalPods *simontype.TargetPodList,
) (score int64, gpuId string) {
	currentFrag := utils.NodeGpuShareFragAmountScore(nodeRes, *typicalPods)

	// partial GPU request
	if podRes.GpuNumber == 1 && podRes.MilliGpu < gpushareutils.MILLI {
		bestCost := math.MaxFloat64
		bestFrag := math.MaxFloat64
		bestGpuId := ""

		for i := 0; i < len(nodeRes.MilliGpuLeftList); i++ {
			if nodeRes.MilliGpuLeftList[i] < podRes.MilliGpu {
				continue
			}

			newNodeRes := nodeRes.Copy()
			newNodeRes.MilliCpuLeft -= podRes.MilliCpu
			newNodeRes.MilliGpuLeftList[i] -= podRes.MilliGpu

			postFrag := utils.NodeGpuShareFragAmountScore(newNodeRes, *typicalPods)
			cost := evaluateNodeFragmentCost(newNodeRes, postFrag, *typicalPods, i, podRes.MilliGpu)

			if bestGpuId == "" || cost < bestCost || (almostEqual(cost, bestCost) && postFrag < bestFrag) {
				bestCost = cost
				bestFrag = postFrag
				bestGpuId = strconv.Itoa(i)
			}
		}

		if bestGpuId == "" {
			return framework.MinNodeScore, ""
		}

		// score: prioritize lower post-placement fragmentation, then larger reduction
		// convert cost minimization to score maximization
		fragReduction := currentFrag - bestFrag
		score = composeScoreFromFrag(currentFrag, bestFrag, fragReduction)
		return score, bestGpuId
	}

	// exclusive or multi-GPU request
	bestCost := math.MaxFloat64
	bestFrag := math.MaxFloat64
	bestGpuId := ""

	candidates := enumerateFeasibleGpuCombinations(nodeRes, podRes)
	for _, ids := range candidates {
		newNodeRes, ok := simulateGpuPlacement(nodeRes, podRes, ids)
		if !ok {
			continue
		}

		postFrag := utils.NodeGpuShareFragAmountScore(newNodeRes, *typicalPods)
		cost := evaluateNodeFragmentCost(newNodeRes, postFrag, *typicalPods, -1, 0)

		if bestGpuId == "" || cost < bestCost || (almostEqual(cost, bestCost) && postFrag < bestFrag) {
			bestCost = cost
			bestFrag = postFrag
			bestGpuId = joinGpuIDs(ids)
		}
	}

	if bestGpuId == "" {
		return framework.MinNodeScore, ""
	}

	fragReduction := currentFrag - bestFrag
	score = composeScoreFromFrag(currentFrag, bestFrag, fragReduction)
	return score, bestGpuId
}

func allocateGpuIdBasedOnFGDPredictorScore(
	nodeRes simontype.NodeResource,
	podRes simontype.PodResource,
	_ simontype.GpuPluginCfg,
	typicalPods *simontype.TargetPodList,
) string {
	_, gpuId := calculateOptimizedScore(nodeRes, podRes, typicalPods)
	return gpuId
}

// -------------------------
// Cost model
// -------------------------

func evaluateNodeFragmentCost(nodeRes simontype.NodeResource, postFrag float64, typicalPods simontype.TargetPodList, chosenGPU int, reqMilli int64) float64 {
	// Main objective: minimize post-placement fragmentation
	const wPostFrag = 1.0

	// Strongly discourage tiny unusable tails
	const wTinyTail = 0.35

	// Prefer packing onto already-used GPUs, preserving fully free GPUs
	const wFreeGpuPreserve = 0.25

	// Penalize residual pieces that are too small for typical pod GPU demands
	const wUnmatch = 0.30

	// Slight penalty for spreading residual capacity across too many half-used GPUs
	const wSpread = 0.15

	tinyTailPenalty := calcTinyTailPenalty(nodeRes)
	freeGpuPreservePenalty := calcFreeGpuPreservePenalty(nodeRes)
	unmatchPenalty := calcUnmatchTypicalPenalty(nodeRes, typicalPods)
	spreadPenalty := calcSpreadPenalty(nodeRes)

	return wPostFrag*postFrag +
		wTinyTail*tinyTailPenalty +
		wFreeGpuPreserve*freeGpuPreservePenalty +
		wUnmatch*unmatchPenalty +
		wSpread*spreadPenalty
}

func calcTinyTailPenalty(nodeRes simontype.NodeResource) float64 {
	// Penalize very small leftovers, e.g. 1~200m, which are often hard to use.
	var penalty float64
	for _, left := range nodeRes.MilliGpuLeftList {
		if left == 0 || left == gpushareutils.MILLI {
			continue
		}
		if left <= 200 {
			penalty += 3.0
		} else if left <= 400 {
			penalty += 1.5
		} else if left <= 600 {
			penalty += 0.5
		}
	}
	return penalty
}

func calcFreeGpuPreservePenalty(nodeRes simontype.NodeResource) float64 {
	// The fewer fully free GPUs, the higher the penalty.
	// Since we minimize cost, this indirectly prefers preserving full GPUs.
	var partialCnt, fullFreeCnt int
	for _, left := range nodeRes.MilliGpuLeftList {
		if left == gpushareutils.MILLI {
			fullFreeCnt++
		} else if left > 0 {
			partialCnt++
		}
	}
	// More partial GPUs means more fragmentation of shape.
	return float64(partialCnt) - 0.3*float64(fullFreeCnt)
}

func calcUnmatchTypicalPenalty(nodeRes simontype.NodeResource, typicalPods simontype.TargetPodList) float64 {
	// Penalize leftovers that cannot fit any typical partial-GPU pod.
	typicalMilliSet := collectTypicalPartialGpuDemands(typicalPods)

	if len(typicalMilliSet) == 0 {
		return 0
	}

	var penalty float64
	for _, left := range nodeRes.MilliGpuLeftList {
		if left == 0 || left == gpushareutils.MILLI {
			continue
		}
		match := false
		for _, req := range typicalMilliSet {
			if left >= req {
				match = true
				break
			}
		}
		if !match {
			penalty += 2.0
		}
	}
	return penalty
}

func calcSpreadPenalty(nodeRes simontype.NodeResource) float64 {
	// Penalize having many partially occupied GPUs.
	var partialCnt int
	for _, left := range nodeRes.MilliGpuLeftList {
		if left > 0 && left < gpushareutils.MILLI {
			partialCnt++
		}
	}
	return float64(partialCnt)
}

func collectTypicalPartialGpuDemands(typicalPods simontype.TargetPodList) []int64 {
	reqMap := map[int64]struct{}{}
	for _, t := range typicalPods {
		tp := t.TargetPodResource
		if tp.GpuNumber == 1 && tp.MilliGpu > 0 && tp.MilliGpu < gpushareutils.MILLI {
			reqMap[tp.MilliGpu] = struct{}{}
		}
	}
	var reqs []int64
	for req := range reqMap {
		reqs = append(reqs, req)
	}
	sort.Slice(reqs, func(i, j int) bool { return reqs[i] < reqs[j] })
	return reqs
}

// -------------------------
// Scoring
// -------------------------

func composeScoreFromFrag(currentFrag, postFrag, reduction float64) int64 {
	// Primary goal: lower postFrag => higher score
	// Secondary goal: larger reduction => slightly higher score
	//
	// We avoid strong sigmoid compression to preserve ranking resolution.
	// Convert to [0, MaxNodeScore].
	const alpha = 0.8 // prioritize post-fragmentation minimum
	const beta = 0.2  // secondary benefit from reduction

	postComponent := 1.0 / (1.0 + postFrag/1000.0)
	reduceComponent := 0.0
	if reduction > 0 {
		reduceComponent = reduction / (reduction + 1000.0)
	}

	raw := alpha*postComponent + beta*reduceComponent
	if raw < 0 {
		raw = 0
	}
	if raw > 1 {
		raw = 1
	}
	return int64(raw * float64(framework.MaxNodeScore))
}

// -------------------------
// Multi-GPU placement enumeration
// -------------------------

func enumerateFeasibleGpuCombinations(nodeRes simontype.NodeResource, podRes simontype.PodResource) [][]int {
	gpuNeed := int(podRes.GpuNumber)
	if gpuNeed <= 0 || gpuNeed > len(nodeRes.MilliGpuLeftList) {
		return nil
	}

	var candidates []int
	for i, left := range nodeRes.MilliGpuLeftList {
		// exclusive allocation usually needs a fully free GPU
		if left == gpushareutils.MILLI {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) < gpuNeed {
		return nil
	}

	var result [][]int
	var path []int
	var dfs func(start int)
	dfs = func(start int) {
		if len(path) == gpuNeed {
			comb := make([]int, len(path))
			copy(comb, path)
			result = append(result, comb)
			return
		}
		for i := start; i < len(candidates); i++ {
			path = append(path, candidates[i])
			dfs(i + 1)
			path = path[:len(path)-1]
		}
	}
	dfs(0)
	return result
}

func simulateGpuPlacement(nodeRes simontype.NodeResource, podRes simontype.PodResource, ids []int) (simontype.NodeResource, bool) {
	newNodeRes := nodeRes.Copy()
	newNodeRes.MilliCpuLeft -= podRes.MilliCpu
	if newNodeRes.MilliCpuLeft < 0 {
		return simontype.NodeResource{}, false
	}

	// For exclusive / multi-GPU requests, consume full GPUs.
	for _, id := range ids {
		if id < 0 || id >= len(newNodeRes.MilliGpuLeftList) {
			return simontype.NodeResource{}, false
		}
		if newNodeRes.MilliGpuLeftList[id] < gpushareutils.MILLI {
			return simontype.NodeResource{}, false
		}
		newNodeRes.MilliGpuLeftList[id] = 0
	}
	return newNodeRes, true
}

func joinGpuIDs(ids []int) string {
	if len(ids) == 0 {
		return ""
	}

	var gpuId string
	for _, id := range ids {
		if gpuId == "" {
			gpuId = strconv.Itoa(id)
		} else {
			gpuId += fmt.Sprintf("%s%d", gpushareutils.DevIdSep, id)
		}
	}
	return gpuId
}

func almostEqual(a, b float64) bool {
	const eps = 1e-9
	return math.Abs(a-b) <= eps
}
