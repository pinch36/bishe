package plugin

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"

	"k8s.io/klog/v2"

	simontype "github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/type"
	gpushareutils "github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/type/open-gpu-share/utils"
	"github.com/hkust-adsl/kubernetes-scheduler-simulator/pkg/utils"
)

// PodGpuPredictor 集成未来得分计算的GPU预测器
type PodGpuPredictor struct {
	// 固定窗口核心字段
	windowMu      sync.Mutex              // 窗口并发安全锁
	windowSize    int                     // 窗口固定大小
	historyWindow []simontype.PodResource // 运行时请求窗口
	// 原有字段
	alpha float64 // 指数平滑系数
	// 新增：未来得分计算依赖参数
	typicalPods      *simontype.TargetPodList // 典型Pod列表
	maxFragThreshold int64                    // 碎片归一化阈值
}

// NewPodGpuPredictor 初始化预测器（新增依赖参数）
func NewPodGpuPredictor(
	windowSize int,
	alpha float64,
	typicalPods *simontype.TargetPodList,
	maxFragThreshold int64,
) (*PodGpuPredictor, error) {
	// 1. 参数校验
	if windowSize <= 0 {
		return nil, fmt.Errorf("windowSize must be >0 (got %d)", windowSize)
	}
	if alpha <= 0 || alpha >= 1 {
		return nil, fmt.Errorf("alpha must be in (0,1) (got %.2f)", alpha)
	}
	if typicalPods == nil {
		klog.V(4).Info("typicalPods is nil, use default list")
		typicalPods = utils.GenerateDefaultTypicalPods() // 兜底默认列表
	}
	if maxFragThreshold <= 0 {
		maxFragThreshold = 10000 // 默认阈值
		klog.V(4).Infof("maxFragThreshold is invalid, use default %d", maxFragThreshold)
	}

	// 2. 初始化空窗口
	historyWindow := make([]simontype.PodResource, 0, windowSize)
	klog.V(4).Infof("Init predictor: window=%d, alpha=%.2f, maxFragThreshold=%d",
		windowSize, alpha, maxFragThreshold)

	// 3. 初始化预测器
	return &PodGpuPredictor{
		windowSize:       windowSize,
		historyWindow:    historyWindow,
		alpha:            alpha,
		typicalPods:      typicalPods,
		maxFragThreshold: maxFragThreshold,
	}, nil
}

// -------------------------- 原有方法保留（无修改） --------------------------
func (p *PodGpuPredictor) PredictNextWindow() (
	milliGpuDist map[int64]float64,
	typicalPodUpdateProb float64,
	reqRate float64,
) {
	milliGpuDist = make(map[int64]float64)
	typicalPodUpdateProb = 0.1
	reqRate = 0.0

	p.windowMu.Lock()
	defer p.windowMu.Unlock()

	if len(p.historyWindow) == 0 {
		klog.V(6).Info("History window is empty, return default milliGPU distribution")
		milliGpuDist[200] = 0.4
		milliGpuDist[300] = 0.3
		milliGpuDist[500] = 0.3
		return
	}

	total := float64(len(p.historyWindow))
	prevDist := make(map[int64]float64)
	for _, podRes := range p.historyWindow {
		if podRes.MilliGpu <= 0 {
			continue
		}
		prevDist[podRes.MilliGpu] += 1.0 / total
	}

	for milli, prob := range prevDist {
		milliGpuDist[milli] = p.alpha*prob + (1-p.alpha)*milliGpuDist[milli]
		if milliGpuDist[milli] < 0 {
			milliGpuDist[milli] = 0
		}
	}

	reqRate = float64(len(p.historyWindow)) / float64(p.windowSize)
	return
}

func (p *PodGpuPredictor) UpdateWindow(podRes simontype.PodResource) {
	if podRes.MilliGpu <= 0 || podRes.GpuNumber <= 0 {
		klog.V(6).Info("Skip invalid pod resource")
		return
	}

	p.windowMu.Lock()
	defer p.windowMu.Unlock()

	p.historyWindow = append(p.historyWindow, podRes)
	if len(p.historyWindow) > p.windowSize {
		p.historyWindow = p.historyWindow[len(p.historyWindow)-p.windowSize:]
	}

	klog.V(6).Infof("Updated history window: current size=%d (max=%d)",
		len(p.historyWindow), p.windowSize)
}

// -------------------------- 新增：未来碎片得分计算（核心） --------------------------
// CalculateFutureFragScore 计算未来碎片预测得分（原插件中的calculateFutureFragScore逻辑）
func (p *PodGpuPredictor) CalculateFutureFragScore(
	nodeRes simontype.NodeResource,
	podRes simontype.PodResource,
) (int64, error) {
	// 1. 预测未来Pod请求特征
	milliGpuDist, typicalPodUpdateProb, _ := p.PredictNextWindow()

	// 2. 复制节点资源（避免修改原数据）
	newNodeRes := nodeRes.Copy()

	// 3. 模拟调度当前Pod（资源扣减安全校验）
	if podRes.GpuNumber == 1 && podRes.MilliGpu < gpushareutils.MILLI { // 共享模式
		// 校验GPU ID合法性
		_, currentGpuId := plugin.CalculateGpuShareFragExtendScore(nodeRes, podRes, p.typicalPods)
		gpuIdx, err := strconv.Atoi(currentGpuId)
		if err != nil {
			return 0, fmt.Errorf("convert GPU ID %s to int failed: %w", currentGpuId, err)
		}
		// 校验资源充足性
		if gpuIdx < 0 || gpuIdx >= len(newNodeRes.MilliGpuLeftList) {
			return 0, fmt.Errorf("GPU index %d out of range (0~%d)", gpuIdx, len(newNodeRes.MilliGpuLeftList)-1)
		}
		if newNodeRes.MilliGpuLeftList[gpuIdx] < podRes.MilliGpu {
			return 0, fmt.Errorf("GPU %d insufficient resource: left=%d, required=%d",
				gpuIdx, newNodeRes.MilliGpuLeftList[gpuIdx], podRes.MilliGpu)
		}
		// 安全扣减资源
		newNodeRes.MilliGpuLeftList[gpuIdx] -= podRes.MilliGpu
	} else { // 独占模式
		// 校验资源充足性
		if newNodeRes.GetTotalMilliGpuLeft() < podRes.MilliGpu*int64(podRes.GpuNumber) {
			return 0, fmt.Errorf("node %s insufficient GPU resource: left=%d, required=%d",
				newNodeRes.NodeName, newNodeRes.GetTotalMilliGpuLeft(), podRes.MilliGpu*int64(podRes.GpuNumber))
		}
		// 安全扣减资源
		newNodeRes, _ = nodeRes.Sub(podRes)
	}

	// 4. 模拟未来Pod请求调度
	for milliGpu, prob := range milliGpuDist {
		if milliGpu <= 0 || prob <= 0 {
			continue
		}
		if rand.Float64() < prob {
			if newNodeRes.GetTotalMilliGpuLeft() < milliGpu {
				klog.V(6).Infof("Node %s no enough resource for future pod (milliGPU=%d), skip",
					newNodeRes.NodeName, milliGpu)
				continue
			}
			futurePodRes := simontype.PodResource{
				MilliGpu:  milliGpu,
				GpuNumber: 1,
			}
			_, futureGpuId := plugin.CalculateGpuShareFragExtendScore(newNodeRes, futurePodRes, p.typicalPods)
			futureGpuIdx, err := strconv.Atoi(futureGpuId)
			if err != nil {
				klog.V(4).Infof("Convert future GPU ID %s to int failed: %v", futureGpuId, err)
				continue
			}
			if futureGpuIdx >= 0 && futureGpuIdx < len(newNodeRes.MilliGpuLeftList) && newNodeRes.MilliGpuLeftList[futureGpuIdx] >= milliGpu {
				newNodeRes.MilliGpuLeftList[futureGpuIdx] -= milliGpu
			}
		}
	}

	// 5. 处理典型Pod列表更新
	futureTypicalPods := p.typicalPods
	if typicalPodUpdateProb > 0 && rand.Float64() < typicalPodUpdateProb {
		if futureTypicalPods != nil {
			futureTypicalPods = utils.UpdateTypicalPods(futureTypicalPods)
		} else {
			klog.V(4).Warning("TypicalPods is nil, skip update")
		}
	}

	// 6. 计算未来碎片值并归一化
	var futureFragValue int64
	if futureTypicalPods != nil {
		futureFragValue = utils.NodeGpuShareFragAmountScore(newNodeRes, *futureTypicalPods)
	} else {
		futureFragValue = utils.NodeGpuShareFragAmountScore(newNodeRes, simontype.TargetPodList{})
	}
	// 归一化：碎片越少 → 得分越高
	normalizedScore := 1.0 - math.Min(1.0, float64(futureFragValue)/float64(p.maxFragThreshold))
	futureScore := int64(math.Round(normalizedScore * float64(100))) // 0~100分

	klog.V(6).Infof("Node %s future frag score: fragValue=%d, normalized=%d",
		nodeRes.NodeName, futureFragValue, futureScore)

	return futureScore, nil
}

// 只读接口（供插件调用）
func (p *PodGpuPredictor) GetWindowSize() int {
	return p.windowSize
}

func (p *PodGpuPredictor) GetCurrentWindowLength() int {
	p.windowMu.Lock()
	defer p.windowMu.Unlock()
	return len(p.historyWindow)
}
