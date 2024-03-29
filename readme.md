
# Snake Game w/ Deep Q Learning
![The snake is alive!](/snakeai.gif)

After recently getting hooked on an original pacman arcade, I became interested in learning to create some simple games of my own. True to my roots, I became simultaneously interested in brushing up on machine learning principles and continuing to flex my Golang muscles. Creating any machine-learning app in Golang as a beginner is an exercise in near futility, as the vast majority of books and examples to learn from are based on Python. Never the less, I persevered, and I'd like to make this code available to any would-be golang gaming and ai enthusiasts. 

## Results
Right now the neural net trains on 5000 games, and has acheived a max score of 40 points. I've had to experiment with tuning the input state, and the reward function to get these results. So far, it doesn't seem like training on more games adds any value.

## Next steps
- [x] Prove that neural net actually learns to play the game
- [x] Help snake avoid infinite loops around the board
- [ ] Persist the derived weights from training so the neural net doesn't have to train every time we start the project

- [ ] Figure out a way to train the snake not to coil up on itself when it's body is in between it's head and the food. I'm guessing I will either need to update the state to represent this, or I will need to tweak the reward function

- [ ] Generalize the code to be useful for more games and pet projects to have fun and explore!

## Resources I learned from for this project
[Alex Pliutau's primer on ebiten engine and snake](https://pliutau.com/ebiten-snake/)

[Mauro Comi's helpful post on a python version of snake and reinforcement learning](https://towardsdatascience.com/how-to-teach-an-ai-to-play-games-deep-reinforcement-learning-28f9b920440a)

[Hand's on Deep Learning with Go](https://github.com/PacktPublishing/Hands-On-Deep-Learning-with-Go)

## Important warnings
Unfortunately, Hands-on Deep Learning with Go is riddled with bugs and outright errors when it comes to their example projects in the book. Chapter 7 about deep Q learning was wholly incomplete and quite a number of key dots in the algorithm were left unconnected. It was only through learning the mathematics behind deep Q learning and the bellman equation, and comparing notes with examples in python that I was able to bridge the gap and successfully complete this project.

This was a toy project I put together over the course of a couple weeks. I was learning many new concepts all at once, and I was not focused on Golang best practices. For that reason, this code could use a lot of improvement. If I have the time, I might get around it it... oh wait, I hear the baby crying, brb.
