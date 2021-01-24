#!/bin/bash
nohup ./berriot --config=testnet/config1.json mine >> logs/1.log &
nohup ./berriot --config=testnet/config2.json mine >> logs/2.log &
nohup ./berriot --config=testnet/config3.json mine >> logs/3.log &
nohup ./berriot --config=testnet/config4.json mine >> logs/4.log &
nohup ./berriot --config=testnet/config5.json mine >> logs/5.log &
