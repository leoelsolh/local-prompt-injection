import glob
import json

from config import (VERBOSE)

files = glob.glob('results/scan_*.json')

if VERBOSE:
    print(f'found {len(files)} files') 

count = {}
total = {}

for path in files: 
    with open(path) as f:
        data = json.load(f)

        for model in data:
            record = data[model]

            for payload in record:
                fields = record[payload] 
                label = fields['label']

                if VERBOSE:
                    print(f"\n{model} | {payload}:\n {fields['label']}")
                
                key = (model, payload, label)
                total_key = (model, payload)
                if key in count:
                    count[key] += 1
                else:
                    count[key] = 1

                if total_key in total:
                    total[total_key] += 1 
                else: 
                    total[total_key] = 1
                    

for key, value in sorted(count.items()):
    group_total = total[key[:2]]
    print(f"{key}: {value}/{group_total}")
