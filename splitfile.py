import os

path = 'D:/X/mimikatz/x64/'
exe_file = path + 'file.exe'
part1 = path + 'firstHalf.exe' 
part2 = path + 'secondHalf.exe'

# Get file size  
exe_size = os.path.getsize(exe_file)
midpoint = exe_size // 2

with open(exe_file, 'rb') as f:
    binary_content = f.read()

with open(part1, 'wb') as f1:
    f1.write(binary_content[:midpoint])

with open(part2, 'wb') as f2:
    f2.write(binary_content[midpoint:])

print(f'{exe_file} split into: ')
print(f'- {part1} ({os.path.getsize(part1)} bytes)')
print(f'- {part2} ({os.path.getsize(part2)} bytes)')
