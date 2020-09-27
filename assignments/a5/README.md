## Homework 5

1. What is TensorFlow? Which company is the leading contributor to TensorFlow?
> TensorFlow is an open source machine learning framework developed by Google.

1. What is TensorRT? How is it different from TensorFlow?
> TensorRT is an SDK developed by Nvidia to leverage GPU processing power with the claim of improving speed during inference by up to 40x. TensorRT is used to optimize performance when running on a GPU and can integrate with other frameworks such as TensorFlow.

1. What is ImageNet? How many images does it contain? How many classes?
> ImageNet is an image database organized according to the WordNet hierarchy.
>- High level categories: 27
>- Total number of non-empty synsets: 21,841
>- Total number of images: 14,197,122

1. Please research and explain the differences between MobileNet and GoogleNet (Inception) architectures.
> MobileNet uses depth-wise separable convolutions to build light weight deep neural networks. This design is well suited for lower powered edge or mobile devices. GoogleNet is an inception architecture that uses 22 convolution layers.

1. In your own words, what is a bottleneck?
> A bottleneck is a layer in a neural network that has less neurons above and below it. This layer is used to compress the representation of the input to reduce the dimensionality of the input.

1. How is a bottleneck different from the concept of layer freezing?
> Layer freezing simply freezes the weights of a pre-trained model while a bottleneck is used to reduce dimensionality. 

1. In the TF1 lab, you trained the last layer (all the previous layers retain their already-trained state). Explain how the lab used the previous layers (where did they come from? how were they used in the process?)
> The previous layers were trained from the ImageNet dataset and were downloaded directly. The weights of the previous layers were used for the final training in the lab.

1. How does a low --learning_rate (step 7 of TF1) value (like 0.005) affect the precision? How much longer does training take?
> Training takes longer but the percision improves a bit.

1. How about a --learning_rate (step 7 of TF1) of 1.0? Is the precision still good enough to produce a usable graph?
> A high learning rate results in a lot of oscillation of the accuracy numbers. Eventually the accuracy stabilized and seemed to produce a usable graph.

1. For step 8, you can use any images you like. Pictures of food, people, or animals work well. You can even use ImageNet images. How accurate was your model? Were you able to train it using a few images, or did you need a lot?
> I trained on images of cats and dogs and used less than 100 images each. The model's accuracy rose to 100% almost immediately and successfully labeled some sample images.

1. Run the TF1 script on the CPU (see instructions above) How does the training time compare to the default network training (section 4)? Why?
> Training with the CPU takes longer because TensorRT is optimized to run on the GPU

1. Try the training again, but this time do export ARCHITECTURE="inception_v3" Are CPU and GPU training times different?
> CPU performance was better

1. Given the hints under the notes section, if we trained Inception_v3, what do we need to pass to replace ??? below to the label_image script? Can we also glean the answer from examining TensorBoard?
```
python -m scripts.label_image --input_layer=??? --input_height=??? --input_width=???  --graph=tf_files/retrained_graph.pb --image=tf_files/flower_photos/daisy/21652746_cc379e0eea_m.jpg
```
> We need to pass 299
