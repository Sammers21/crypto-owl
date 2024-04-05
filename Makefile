default: tether usdc uniswap weth9

tether:
	cd abis/tether && abigen --abi=TetherToken.abi --pkg=tether --out=TetherToken.go --alias _totalSupply=_totalSupplyRenamed
.PHONY: tether

usdc:
	cd abis/FiatTokenProxy && abigen --abi=FiatTokenProxy.abi --pkg=FiatTokenProxy --out=FiatTokenProxy.go
	cd abis/FiatTokenV22 && abigen --abi=FiatTokenV22.abi --pkg=FiatTokenV22 --out=FiatTokenV22.go
.PHONY: usdc

weth9:
	cd abis/WETH9 && abigen --abi=WETH9.abi --pkg=WETH9 --out=WETH9.go

uniswap:
	cd abis/uniswap && abigen --abi=SwapRouter02.abi --pkg=uniswap --out=SwapRouter02.go

.PHONY: usdc