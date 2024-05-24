import functools
import matplotlib.pyplot as plt
import numpy as np
import matplotlib.animation as animation

HIST_BINS = [2**i for i in range(1, 14)]
lab = ("2", "4", "8", "16", "32", "64", "128", "256", "512", "1024", "2048", "4096")

data_all = np.genfromtxt("bar_plot/highest_scores.csv", delimiter=",")
data = data_all[0]

def animate(frame_number, bar_container, text):
    # Simulate new data coming in.
    data = data_all[frame_number]
    n, _ = np.histogram(data, HIST_BINS)
    for count, rect in zip(n, bar_container.patches):
        rect.set_height(count)
        
    text.set_text(f'Generation: {frame_number+1}')
    
    return bar_container.patches + [text]

n, _ = np.histogram(data, HIST_BINS)
fig, ax = plt.subplots()
bar = ax.bar(lab[:len(n)], n)

# Add text annotation
text = ax.text(0.02, 0.95, '', transform=ax.transAxes, fontsize=12, verticalalignment='top')

ax.set_title("Highest achieved number by bot")
ax.set_ylim(top=400)
anim = functools.partial(animate, bar_container=bar, text=text)
ani = animation.FuncAnimation(fig, anim, frames=2048, interval=50, repeat=False, blit=True)
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