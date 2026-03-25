package plugin

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	simontype "github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/type"
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
			fmt.Sprintf("failed to get nodeRes(%s)\n", nodeName),
		)
	}
	nodeRes := *nodeResPtr

	podRes := utils.GetPodResource(pod)
	if !utils.IsNodeAccessibleToPod(nodeRes, podRes) {
		return framework.MinNodeScore, framework.NewStatus(
			framework.Error,
			fmt.Sprintf("Node (%s) %s does not match GPU type request of pod %s\n",
				nodeName, nodeRes.Repr(), podRes.Repr()),
		)
	}

	score, _ := calculateGpuShareFragExtendScore(nodeRes, podRes, p.typicalPods)
	return score, framework.NewStatus(framework.Success)
}

func (p *FGDPredictorScorePlugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

func allocateGpuIdBasedOnFGDPredictorScore(
	nodeRes simontype.NodeResource,
	podRes simontype.PodResource,
	_ simontype.GpuPluginCfg,
	typicalPods *simontype.TargetPodList,
) string {
	_, gpuId := calculateGpuShareFragExtendScore(nodeRes, podRes, typicalPods)
	return gpuId
}
