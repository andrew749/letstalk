import os

def template_file(file_path: str, out_path: str, fill_in_dict={}):
    with open(file_path, 'r') as input_file:
        with open(out_path, 'w') as out_file:
            data = input_file.read()
            formatted_data = data.format(fill_in_dict)
            out_file.write(formatted_data)
