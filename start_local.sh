#!/bin/bash
nohup ./berriot --config=local/config1.json mine >> logs/1.log &
nohup ./berriot --config=local/config2.json mine >> logs/2.log &
nohup ./berriot --config=local/config3.json mine >> logs/3.log &
nohup ./berriot --config=local/config4.json mine >> logs/4.log &
nohup ./berriot --config=local/config5.json mine >> logs/5.log &
