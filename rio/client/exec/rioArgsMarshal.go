package rioclient

import (
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
)

func UnpackArgs(
	wareID api.WareID,
	path string,
	filt api.FilesetUnpackFilter,
	placementMode rio.PlacementMode,
	warehouses []api.WarehouseLocation,
	monitor rio.Monitor,
) ([]string, error) {
	// Required args.
	args := []string{"unpack", "--format=json"}

	// Append filters if specified.
	//  (We could just pass 'em all even when emptystr, but let's be nice to readers of 'ps'.)
	filtStr := filt.String()
	if filtStr != "" {
		args = append(args, "--filter", filtStr)
	}

	// Append placement mode if specified.
	if placementMode != "" {
		args = append(args, "--placer="+string(placementMode))
	}

	// Append warehouses.
	//  Giving this argument repeatedly forms a list in the rio CLI.
	for _, wh := range warehouses {
		args = append(args, "--source="+string(wh))
	}

	// Append monitor options.
	//  (Of which there are currently none meaningful implemented.)

	// Suffix the main bits.
	//  This is last so we can use the "--" to terminate acceptance of flags
	//  (which is important because, well, what if someone really does want
	//  to unpack into path "--lol"?).
	args = append(args, "--", wareID.String(), path)

	// Done!
	return args, nil
}

func PackArgs(
	packType api.PackType,
	path string,
	filt api.FilesetPackFilter,
	warehouse api.WarehouseLocation,
	monitor rio.Monitor,
) ([]string, error) {
	// Required args.
	args := []string{"pack", "--format=json"}

	// Append filters if specified.
	//  (We could just pass 'em all even when emptystr, but let's be nice to readers of 'ps'.)
	filtStr := filt.String()
	if filtStr != "" {
		args = append(args, "--filter", filtStr)
	}

	// Append warehouse.
	if warehouse != "" {
		args = append(args, "--target="+string(warehouse))
	}

	// Append monitor options.
	//  (Of which there are currently none meaningful implemented.)

	// Suffix the main bits.
	//  This is last so we can use the "--" to terminate acceptance of flags
	//  (which is important because, well, what if someone really does want
	//  to unpack into path "--lol"?).
	args = append(args, "--", string(packType), path)

	// Done!
	return args, nil
}
