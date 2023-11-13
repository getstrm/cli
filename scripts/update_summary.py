import re
from glob import glob

text = ""
with open('./pace/docs/SUMMARY.md', 'r') as file:
  text = file.read()
with open('./pace/docs/SUMMARY.md', 'w') as file:
  replacement = ""
  for fn in sorted([f[len('./generated_docs/'):] for f in glob('./generated_docs/pace*')]):
    filename = fn.replace('_', "\\_")
    replacement+=f"""{'  '*len(fn.split('_'))}* [{fn.split('_')[-1].replace('.md','')}](cli-docs/{filename})\n"""
  out = re.sub('([\s\S]*\[CLI Reference\]\(reference/cli-reference/README.md\))[\S\s]*?\n(\*[\s\S]*)', f'\\1\n{replacement}\\2', text)
  file.write(out)
