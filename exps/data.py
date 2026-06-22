""" defines util function for interacting with the experiment data """

import os
import pandas as pd
from time import localtime, strftime


def create_output_folder():
    cwd = os.getcwd()
    time_string = strftime("%Y-%m-%d-%H-%M-%S", localtime())
    output_folder = cwd + '/results/results-' + time_string + '/'

    os.makedirs(output_folder)
    os.makedirs(output_folder + 'sim_out/')
    os.makedirs(output_folder + 'gold_out/')

    return output_folder


def data_columns():
    return ['benchmark', 'env', 'param1', 'param2', 'time']


def initialize_gold_data():
    data = pd.DataFrame(columns=data_columns())
    return data


def collect_gold_data(exps):
    gold = pd.DataFrame(columns=data_columns())

    for exp in exps:
        gold = gold.append(exp.gold())

    return gold


def gold_file_name(output_folder, exp_name):
    path = output_folder + '/' + exp_name + '_gold.csv'
    return path


def save_gold_data(data, output_folder, exp_name):
    path = gold_file_name(output_folder, exp_name)
    data.to_csv(path)
    return path
