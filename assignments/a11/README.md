## Results

video link: https://www.dropbox.com/sh/6k9fhwi58q1vt2v/AAAKlmqdNLFZsrV74dgasz70a?dl=0

#### Model 1

- Parameters:
    - density_first_layer = 64
    - density_second_layer = 32
    - batch_size = 128

|Average Reward|Total Above 200|
|--------------|---------------|
|203.47|56|

As this is the first attempted model, I did not want to modify too much from the baseline configuration. I increased the density of the layers in order to prevent under fitting of the model and increased batch size to improve the accuracy of the estimator when calculating the weights.

#### Model 2

- Parameters:
    - density_first_layer = 64
    - density_second_layer = 32
    - batch_size = 128
    - epsilon_decay = 0.99
    - learning rate = 0.005

This model I updated the learning rate and epsilon decay. The goal here was to see if training speed could be improved while still building an accurate model. By lowering the epsilon decay from 0.995 to 0.99 the model will include less random movements at each step and the the increase in learning rate would result in a model that converges faster. Unfortunatley, these changes were too extreem and after initial improvements the model overshort the gradient descent minima and became worse and worse over time.

#### Model 3

- Parameters:
    - density_first_layer = 128
    - density_second_layer = 64
    - batch_size = 128
    - epsilon_decay = 0.99
    - learning rate = 0.0005
    - epochs = 2


This model the learning rate was dropped by an order of magnitude to see if that would prevent the errors in the previous model. The learning rate and epsilon decay with this model suffered a different problem than the previous one. The slower learning rate and faster epsilon decay caused the model to learn that it was rewarded for "hovering" instead of completing landings. This also significantly increased training time since the game had to time out each time the simulation stayed in a hovering state. This model was ultimately killed before training was completed.

#### model 4

- Parameters:
    - density_first_layer = 256
    - density_second_layer = 128
    - batch_size = 128
    - epsilon_decay = 0.995
    - learning rate = 0.0005
    - epochs = 2

|Average Reward|Total Above 200|
|--------------|---------------|
|236.98|86|

I moved the learning rate to 0.0005 and increased the layer densities to 256 and 128. I set training to terminate after 500 episodes so this model only reached and average of 187 when training completed. That said, this model showed a significant improvement over model 1. It should also be said that while the average and total above seem like good numbers there were still several landings that "killed" our astronauts with large negative scores.

```
Starting Testing of the trained model...
36  : Episode || Reward:  -159.72288036685492
44  : Episode || Reward:  -234.27213660201397
98  : Episode || Reward:  -227.9101467026177
```


#### Model 5

- parameters:
    - density_first_layer = 512
    - density_second_layer = 128
    - batch_size = 128
    - epsilon_decay = 0.99
    - learning rate = 0.001
    - epochs = 2

|Average Reward|Total Above 200|
|--------------|---------------|
|234.19|89|

The goal of this model is to improve the training speed and try not to kill any astronauts. The first layer density was increased to 512 in order to capture more features before passing to the second layer that remained unchanged at 128. In order to improve training speed the learning rate was moved back to 0.001 while the epsilon decay was moved to 0.99. This seems like the ideal spot in regard to training speed since this model reached an average score of 200 in just over 300 episodes. While the average reward was a little less the total above 200 outperformed model 4, and even more important, this model did not kill any astronauts with a low score of 102. This diffence could possibly be explained by the faster epsilon decay causing this model to be more conservative but "safer" that model 4 but model 4 learning how to maximize the highest score more often.

## Questions

1) What parameters did you change, and what values did you use?

> see above

2) Why did you change these parameters?

> see above

3) Did you try any other changes (like adding layers or changing the epsilon value) that made things better or worse?

> see above

4) Did your changes improve or degrade the model? How close did you get to a test run with 100% of the scores above 200?

> See above. None of the models produced test runs with all scores above 200.

5) Based on what you observed, what conclusions can you draw about the different parameters and their values?

> see above

6) What is the purpose of the epsilon value?

> This introduces randomness in the model in order to identify beneficial choices that have not been learned by the previous models.

7) Describe "Q-Learning".

> Q-learning is a learning method that seeks to find the best reward for any given state, in this scase the "reward". Through learning from previous models as well as random input (epsilon) the learning process creates a model that will maximize the total reward based on a given state.
