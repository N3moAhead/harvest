#!/usr/bin/env python3
import sys
import subprocess
import os
import argparse

def get_git_authors(file_path: str) -> list[str]:
    """
    Executes 'git log' for a given file and returns a list of unique author names.
    """
    command = ['git', 'log', '--follow', '--pretty=format:%an', '--', file_path]
    
    try:
        result = subprocess.run(
            command,
            capture_output=True,
            text=True,
            check=True,
            encoding='utf-8'
        )
        
        authors_output = result.stdout.strip()
        if not authors_output:
            return []
            
        authors_list = authors_output.split('\n')
        unique_authors = list(dict.fromkeys(authors_list))
        return unique_authors
        
    except FileNotFoundError:
        print("Error: 'git' command not found. Please ensure Git is installed and in your PATH.")
        sys.exit(1)
    except subprocess.CalledProcessError as e:
        print(f"  -> Warning: Could not get authors for '{file_path}'. Git error: {e.stderr.strip()}")
        return []

def update_author_line_in_file(file_path: str, authors: list[str]):
    """
    Prepends or updates a comment line with author names in the specified file.
    If the line exists and is up-to-date, the file is skipped.
    """
    if not authors:
        print(f"  -> Skipping '{os.path.basename(file_path)}' (no authors found).")
        return

    author_string = ", ".join(authors)
    new_comment_line = f"// Authors: {author_string}\n" # English comment prefix

    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
        
        content_to_write = []
        action = ""

        # Check if an author line already exists at the top
        if lines and lines[0].startswith('// Authors:'): # English check
            # Case 1: The line is identical, nothing to do.
            if lines[0].strip() == new_comment_line.strip():
                print(f"  -> Skipping '{os.path.basename(file_path)}' (authors are already up-to-date).")
                return
            
            # Case 2: The line exists but needs an update.
            content_to_write = [new_comment_line] + lines[1:]
            action = "updated"
        else:
            # Case 3: No author line exists, add a new one.
            content_to_write = [new_comment_line] + lines
            action = "added"

        # Write the new content back to the file
        with open(file_path, 'w', encoding='utf-8') as f:
            f.writelines(content_to_write)
        
        print(f"  -> ✅ Author line {action} in '{os.path.basename(file_path)}'.")

    except Exception as e:
        print(f"  -> Error: Failed to process file '{file_path}': {e}")

def main():
    """
    Main function to parse arguments and process files recursively.
    """
    parser = argparse.ArgumentParser(
        description="Prepends or updates a Git author comment in files. Searches directories recursively.",
        epilog="Example: python update_authors.py ./src -e .js .ts .py"
    )
    
    parser.add_argument(
        "directory",
        help="The directory to search recursively."
    )
    
    parser.add_argument(
        "-e", "--extensions",
        required=True,
        nargs="+",
        help="One or more file extensions to process (e.g., .py .java .ts)."
    )
    
    args = parser.parse_args()
    
    # --- Input Validation ---
    if not os.path.isdir(args.directory):
        print(f"Error: Directory '{args.directory}' not found.")
        sys.exit(1)
        
    if not os.path.isdir('.git'):
        print("Error: This script must be run from within a Git repository.")
        sys.exit(1)
        
    valid_extensions = tuple(ext if ext.startswith('.') else f".{ext}" for ext in args.extensions)
    
    # --- File Discovery ---
    files_to_process = []
    for root, _, files in os.walk(args.directory):
        for file in files:
            if file.endswith(valid_extensions):
                files_to_process.append(os.path.join(root, file))

    if not files_to_process:
        print(f"No files with extensions '{' '.join(valid_extensions)}' found in directory '{args.directory}'.")
        return
        
    # --- Processing ---
    total_files = len(files_to_process)
    print(f"\nFound {total_files} matching file(s). Starting process...")
    
    for i, file_path in enumerate(files_to_process, 1):
        print(f"[{i}/{total_files}] Processing: {file_path}")
        authors = get_git_authors(file_path)
        update_author_line_in_file(file_path, authors)
        
    print("\nProcessing complete.")

if __name__ == "__main__":
    main()