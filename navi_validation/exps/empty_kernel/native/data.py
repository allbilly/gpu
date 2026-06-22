import pandas as pd
from sklearn.linear_model import LinearRegression
import matplotlib.pyplot as plt
from exp import run_benchmark_on_gpu
def run_native(repeat, data):

    for p1 in [2**i for i in range(0,15)]:
        for p2 in range(1,9):
            times = []
            for i in range(0, repeat):
                time = run_native_once(p1,p2)
                times.append(time)
            entry = ["empty kernel","native",p1, p2,sum(times[5:])/5]
            data.loc[len(data)] =entry

def run_native_once(param1, param2):
    return run_benchmark_on_gpu('empty kernel','./kernel {0} {1}','.',param1,param2)

def plot_data(data):
    df = data.pivot(index = 'param1', columns = 'param2', values = 'time')
    linear_regressor = LinearRegression()
    x = df.index.values.reshape(-1, 1)
    for col in df.columns:
        linear_regressor.fit(x,df[col].values.reshape(-1, 1))
        time_pred = linear_regressor.predict(x)
        plt.plot(x,time_pred)
    plt.savefig("empty_kernel.png")

def main():
    data = pd.DataFrame(columns=['benchmark','env','param1','param2','time'])

    run_native(10,data)
    data.to_csv('./result.csv')
    data = pd.read_csv('./result.csv')
    plot_data(data)


if __name__ == "__main__":
    main()
