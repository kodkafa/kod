import sys
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("--name", required=True)
args = parser.parse_args()

print("╔════════════════════════════════════╗")
print(f"║  Hello {args.name: <26} ║")
print("║  from KODKAFA Python Plugin!       ║")
print("╚════════════════════════════════════╝")
