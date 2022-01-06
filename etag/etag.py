#!/usr/bin/env python3

""""
Author: Jay Lux Ferro
Specs:  A Simple E-Tag Generation Technique
"""

class ETag:
    def __init__(self, seed=0):
        # seed -> int
        self.seed = seed

    def generate(self, timestamp):
        # timestamp -> int (time in milliseconds)
        return self.bin2hex(bin(self.seed)[2:] + bin(timestamp)[2:])

    def bin2hex(self, binary_data):
        return hex(int(binary_data, 2))[2:].upper()

    def hex2bin(self, hex_data):
        return bin(int(hex_data, 16))[2:]

    def is_valid(self, e_tag, timestamp, seed=0):
        self.seed = seed
        return True if self.generate(timestamp) == e_tag.upper() else False

