tether:
	cd tether && abigen --abi=TetherToken.abi --pkg=tether --out=TetherToken.go --alias _totalSupply=_totalSupplyRenamed
.PHONY: tether