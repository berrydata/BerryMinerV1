#!/bin/bash
nohup ./berriot --config=test_cfgs/config1.json mine >> logs/1.log &
nohup ./berriot --config=test_cfgs/config2.json mine >> logs/2.log &
nohup ./berriot --config=test_cfgs/config3.json mine >> logs/3.log &
nohup ./berriot --config=test_cfgs/config4.json mine >> logs/4.log &
nohup ./berriot --config=test_cfgs/config5.json mine >> logs/5.log &
