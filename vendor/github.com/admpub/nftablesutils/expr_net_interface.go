package nftablesutils

import "github.com/google/nftables/expr"

// SetIIF equals input-interface
func SetIIF(iface string) Exprs {
	exprs := []expr.Any{
		ExprIIFName(),
		ExprCmpEqIFName(iface),
	}

	return exprs
}

// SetOIF equals output-interface
func SetOIF(iface string) Exprs {
	exprs := []expr.Any{
		ExprOIFName(),
		ExprCmpEqIFName(iface),
	}

	return exprs
}

// SetNIIF not equals input-interface
func SetNIIF(iface string) Exprs {
	exprs := []expr.Any{
		ExprIIFName(),
		ExprCmpNeqIFName(iface),
	}

	return exprs
}

// SetNOIF not equals output-interface
func SetNOIF(iface string) Exprs {
	exprs := []expr.Any{
		ExprOIFName(),
		ExprCmpNeqIFName(iface),
	}

	return exprs
}
