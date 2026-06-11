import glob
import json

from config import (VERBOSE)

files = glob.glob('results/scan_*.json')

if VERBOSE:
    print(f'found {len(files)} files') 

count = {}
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
                if key in count:
                    count[key] += 1
                else:
                    count[key] = 1

for key, value in count.items():
    print(key, value)