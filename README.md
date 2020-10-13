# Factom Chain ID Generator

Allows you to generate chains with an arbitrary prefix, such as the ones used for identity chains. 

### Example

```bash
$ ./factom-chainid-generator -workers=4 -target="c0ffee" -base="Coffee Chain"
Using base: [Coffee Chain]
Using target: c0ffee
Chain found with extid 3877731: c0ffee6b98e3b4391ed7af22f16310ff89b303e4bb5c138aadc797328c250479 (in 1.1229892s)
ExtIDs: "Coffee Chain" "3877731"
```

You can now create this chain via factom-cli:
```bash
$ factom-cli addchain -n "Coffee Chain" -n "3877731" EC...
```