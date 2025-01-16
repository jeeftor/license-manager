#!/usr/bin/env perl
use strict;
use warnings;

sub greet {
    my $name = shift;
    return "Hello, $name!";
}

print greet("World") . "\n";
