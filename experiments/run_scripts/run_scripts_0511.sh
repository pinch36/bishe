#!/bin/bash
# cat run_scripts_0511.sh | while read i; do printf "%q\n" "$i"; done | xargs --max-procs=16 -I CMD bash -c CMD

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_default
EXPDIR="experiments/2023_0511/openb_pod_list_default/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_default -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_cpu050
EXPDIR="experiments/2023_0511/openb_pod_list_cpu050/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_cpu050 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_cpu100
EXPDIR="experiments/2023_0511/openb_pod_list_cpu100/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_cpu100 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_cpu200
EXPDIR="experiments/2023_0511/openb_pod_list_cpu200/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_cpu200 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_cpu250
EXPDIR="experiments/2023_0511/openb_pod_list_cpu250/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_cpu250 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_gpushare100
EXPDIR="experiments/2023_0511/openb_pod_list_gpushare100/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_gpushare100 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_gpushare40
EXPDIR="experiments/2023_0511/openb_pod_list_gpushare40/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_gpushare40 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_gpushare60
EXPDIR="experiments/2023_0511/openb_pod_list_gpushare60/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_gpushare60 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_gpushare80
EXPDIR="experiments/2023_0511/openb_pod_list_gpushare80/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_gpushare80 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_gpuspec10
EXPDIR="experiments/2023_0511/openb_pod_list_gpuspec10/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_gpuspec10 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_gpuspec20
EXPDIR="experiments/2023_0511/openb_pod_list_gpuspec20/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_gpuspec20 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_gpuspec25
EXPDIR="experiments/2023_0511/openb_pod_list_gpuspec25/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_gpuspec25 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_gpuspec33
EXPDIR="experiments/2023_0511/openb_pod_list_gpuspec33/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_gpuspec33 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_multigpu20
EXPDIR="experiments/2023_0511/openb_pod_list_multigpu20/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_multigpu20 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_multigpu30
EXPDIR="experiments/2023_0511/openb_pod_list_multigpu30/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_multigpu30 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_multigpu40
EXPDIR="experiments/2023_0511/openb_pod_list_multigpu40/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_multigpu40 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

# 07, FGDPredictor, FGDPredictor, share, max @ openb_pod_list_multigpu50
EXPDIR="experiments/2023_0511/openb_pod_list_multigpu50/07-FGDPredictor/1.3/42" && mkdir -p ${EXPDIR} && touch "${EXPDIR}/terminal.out" && python3 scripts/generate_config_and_run.py -d "${EXPDIR}" -e -b -f data/openb_pod_list_multigpu50 -FGDPredictor 1000 -gpusel FGDPredictor -dimext share -norm max -tune 1.3 -tuneseed 42 --shuffle-pod=true -z "${EXPDIR}/snapshot/ds01" | tee -a "${EXPDIR}/terminal.out" && python3 scripts/analysis.py -f -g ${EXPDIR} | tee -a "${EXPDIR}/terminal.out" 

