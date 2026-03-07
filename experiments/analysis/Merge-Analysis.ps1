# Merge-Analysis.ps1 - Windows 版分析文件合并脚本
$targetDir = ".\analysis_results"
# 创建结果目录（不存在则创建）
if (-not (Test-Path $targetDir)) { New-Item -ItemType Directory -Path $targetDir | Out-Null }

# 定义需要合并的文件类型（与原脚本一致）
$analysisFiles = @(
    "analysis_allo_discrete.csv",
    "analysis_frag_ratio_discrete.csv",
    "analysis_frag_discrete.csv",
    "analysis_fail.csv"
)

# 遍历 data 目录下所有实验文件夹
Get-ChildItem -Path ..\data -Recurse -File | Where-Object {
    $analysisFiles -contains $_.Name
} | Group-Object -Property Name | ForEach-Object {
    $fileName = $_.Name
    $outputPath = Join-Path $targetDir $fileName
    Write-Host "合并 $fileName 到 $outputPath"

    # 清空目标文件（避免重复）
    if (Test-Path $outputPath) { Clear-Content $outputPath }

    # 逐行合并（保留表头一次，后续跳过表头）
    $isFirstFile = $true
    $_.Group | ForEach-Object {
        $filePath = $_.FullName
        $content = Get-Content $filePath -Encoding UTF8
        if ($isFirstFile) {
            # 第一个文件保留表头
            $content | Add-Content -Path $outputPath -Encoding UTF8
            $isFirstFile = $false
        } else {
            # 后续文件跳过表头（取第2行开始）
            $content | Select-Object -Skip 1 | Add-Content -Path $outputPath -Encoding UTF8
        }
    }
}

Write-Host "合并完成！结果保存在 $targetDir"