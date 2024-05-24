import functools

import matplotlib.pyplot as plt
import numpy as np

import matplotlib.animation as animation

HIST_BINS = [2**i for i in range(1, 12)]

data_all = np.genfromtxt("highest_scores.csv", delimiter=",")
data = data_all[0]

n, _ = np.histogram(data, HIST_BINS)

def animate(frame_number, bar_container):
    # Simulate new data coming in.
    data = data_all[frame_number]
    n, _ = np.histogram(data, HIST_BINS)
    for count, rect in zip(n, bar_container.patches):
        rect.set_height(count)

    return bar_container.patches


labels = HIST_BINS[1:]
fig, ax = plt.subplots()
#_, _, bar_container = ax.hist(data, HIST_BINS, lw=1, ec="yellow", fc="green", alpha=0.5)
bar_container= ax.bar(labels, n)

#ax.set_ylim(top=55)
anim = functools.partial(animate, bar_container=bar_container)
ani = animation.FuncAnimation(fig, anim, 2048, repeat=False, blit=True)
plt.show()

'''
import matplotlib.pyplot as plt
import numpy as np

data = np.genfromtxt("highest_scores.csv", delimiter=",")

d0 = data[0,:]

#plt.bar(range(400), d0)
#plt.show()

plt.hist(d0)
plt.show()
'''