# Baldor Food

This repository runs a local web server that allows you to browse items on
[baldorfood.com](https://www.baldorfood.com/) and add them to your cart.

It is used to get around the slow page loads on the actual site.

## Usage

1. Run `go run cmd/index/main.go`.
  - If you get this error message: `cannot create new index, path already
    exists`, just delete the `product_index` folder and rerun the command.
2. Run `go run cmd/web/main.go`.
3. Visit `localhost:8080` to browse items and add them to your cart.
4. To checkout, use [https://www.baldorfood.com/cart](https://www.baldorfood.com/cart).

## Notes

- Salmon is categorized under `specialtygrocery`, not `meatseafood`.
