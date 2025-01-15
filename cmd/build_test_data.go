// cmd/build_test_data.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildTestDataCmd = &cobra.Command{
	Use:   "build-test-data",
	Short: "Generate test files for all supported languages",
	Long:  `Creates a test_data directory with hello world programs in all supported languages`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildTestData()
	},
}

func init() {
	rootCmd.AddCommand(buildTestDataCmd)
}

// Language specific hello world templates
var templates = map[string][]LanguageTemplate{
	"python": {
		{
			Filename: "hello.py",
			Content: `#!/usr/bin/env python3

def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
`,
		},
	},
	"ruby": {
		{
			Filename: "hello.rb",
			Content: `#!/usr/bin/env ruby

puts "Hello, World!"
`,
		},
	},
	"javascript": {
		{
			Filename: "hello.js",
			Content: `console.log("Hello, World!");
`,
		},
		{
			Filename: "component.jsx",
			Content: `import React from 'react';

const HelloWorld = () => {
    return <div>Hello, World!</div>;
};

export default HelloWorld;
`,
		},
	},
	"typescript": {
		{
			Filename: "hello.ts",
			Content: `function sayHello(): void {
    console.log("Hello, World!");
}

sayHello();
`,
		},
		{
			Filename: "component.tsx",
			Content: `import React from 'react';

const HelloWorld: React.FC = () => {
    return <div>Hello, World!</div>;
};

export default HelloWorld;
`,
		},
	},
	"java": {
		{
			Filename: "HelloWorld.java",
			Content: `public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}
`,
		},
	},
	"go": {
		{
			Filename: "hello.go",
			Content: `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
`,
		},
	},
	"c": {
		{
			Filename: "hello.c",
			Content: `#include <stdio.h>

int main() {
    printf("Hello, World!\n");
    return 0;
}
`,
		},
		{
			Filename: "hello.h",
			Content: `#ifndef HELLO_H
#define HELLO_H

void say_hello(void);

#endif /* HELLO_H */
`,
		},
	},
	"cpp": {
		{
			Filename: "hello.cpp",
			Content: `#include <iostream>

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}
`,
		},
		{
			Filename: "hello.hpp",
			Content: `#ifndef HELLO_HPP
#define HELLO_HPP

namespace hello {
    void say_hello();
}

#endif /* HELLO_HPP */
`,
		},
	},
	"csharp": {
		{
			Filename: "Hello.cs",
			Content: `using System;

class Hello {
    static void Main() {
        Console.WriteLine("Hello, World!");
    }
}
`,
		},
	},
	"php": {
		{
			Filename: "hello.php",
			Content: `<?php
echo "Hello, World!\n";
`,
		},
	},
	"swift": {
		{
			Filename: "hello.swift",
			Content: `print("Hello, World!")
`,
		},
	},
	"rust": {
		{
			Filename: "hello.rs",
			Content: `fn main() {
    println!("Hello, World!");
}
`,
		},
	},
	"shell": {
		{
			Filename: "hello.sh",
			Content: `#!/bin/bash
echo "Hello, World!"
`,
		},
		{
			Filename: "hello.bash",
			Content: `#!/bin/bash
echo "Hello, World!"
`,
		},
	},
	"yaml": {
		{
			Filename: "config.yml",
			Content: `greeting:
  message: "Hello, World!"
  language: "English"
`,
		},
		{
			Filename: "config.yaml",
			Content: `greeting:
  message: "Hello, World!"
  language: "English"
`,
		},
	},
	"perl": {
		{
			Filename: "hello.pl",
			Content: `#!/usr/bin/env perl
use strict;
use warnings;

print "Hello, World!\n";
`,
		},
		{
			Filename: "Hello.pm",
			Content: `package Hello;
use strict;
use warnings;

sub say_hello {
    print "Hello, World!\n";
}

1;
`,
		},
	},
	"r": {
		{
			Filename: "hello.r",
			Content: `print("Hello, World!")
`,
		},
	},
	"html": {
		{
			Filename: "index.html",
			Content: `<!DOCTYPE html>
<html>
<head>
    <title>Hello World</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>
`,
		},
	},
	"xml": {
		{
			Filename: "hello.xml",
			Content: `<?xml version="1.0" encoding="UTF-8"?>
<greeting>
    <message>Hello, World!</message>
</greeting>
`,
		},
	},
	"css": {
		{
			Filename: "style.css",
			Content: `body {
    font-family: sans-serif;
}

.greeting {
    color: blue;
    font-size: 24px;
}
`,
		},
	},
	"scss": {
		{
			Filename: "style.scss",
			Content: `$color: blue;

.greeting {
    color: $color;
    font-size: 24px;
    
    &:hover {
        color: darken($color, 10%);
    }
}
`,
		},
	},
	"sass": {
		{
			Filename: "style.sass",
			Content: `$color: blue

.greeting
    color: $color
    font-size: 24px
    
    &:hover
        color: darken($color, 10%)
`,
		},
	},
	"lua": {
		{
			Filename: "hello.lua",
			Content: `print("Hello, World!")
`,
		},
	},
}

type LanguageTemplate struct {
	Filename string
	Content  string
}

func buildTestData() error {
	baseDir := "test_data"

	// Create base directory
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %v", err)
	}

	// Create language directories and files
	for lang, files := range templates {
		langDir := filepath.Join(baseDir, lang)
		if err := os.MkdirAll(langDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %v", lang, err)
		}

		for _, file := range files {
			filePath := filepath.Join(langDir, file.Filename)
			if err := os.WriteFile(filePath, []byte(file.Content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", filePath, err)
			}
			fmt.Printf("Created %s\n", filePath)
		}
	}

	fmt.Printf("\nSuccessfully created test files in %s directory\n", baseDir)
	return nil
}
