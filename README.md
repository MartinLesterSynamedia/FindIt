# FindIt

## Description
FindIt is a CAPTCHA that is dynamically generated from a known set of images. It asks the simple question "Find these key images in this larger image". Due to the random process the key and background images are constructed it should require significant image processing power to identify the location of the key images. As there should be at least 3 keys it is also very unlikely that you can randomly guess the lcoation of the keys with the search image.

## Process
The process is split into several distinct sections. Though they can be run linearly, the intent is they could each run completely asyncronosly so allowing the process of generating the CAPTCHA to be much more scaleable. 

### Images
The system is based on 2 sets of images, the keys and the backgrounds. The keys should be noticably different to their background. The default example is jungle animals as the keys (parrot, spider, monkey, etc) and the backgrounds are jungle images.

The user is presented with 3 key images and a dynamically generated search image which is constructed from 3 layers. Obfuscator, background & keys.

#### Preprocessing
The first process is to normalise the input images. This converts all images to jpeg and reduces them to a maximum width and height. This saves on disk space and memory when performing other processing. The normalising process only needs to be performed once.

To speed up the process of generating the search image the key and background images should be rotated, scalled, skewed, flipped, etc to generate a collection of altered original images. Each generated image can be used multiple times, but they should also be reproduced on some cron, depending on the rate of CAPTCHA generation and available resources.

#### Alpha Blend
Similar to the Preprocessing a set of Alpha blends must be generated so the selected images can be quickly blended into the final search image.

### Obfuscator
This is the first layer of the search image. It is any randomly generated image using the 3 keys. This is done to prevent the keys being found by simply looking for unique colours within them. There can be many processes to produce this layer e.g. plasma map, partial blurred images, spirals, etc For simplicity the initial implementation will use cropped sections of the keys.

### Background
This uses a random selection of prepocessed background images and combines them with the aplh blends. This is produced as its own image and stored. Similar to the other generated images they are expected to be reproduced on a cron.

### The Search Image
To generate the search image all that is now required is to randomly select 3 keys, a matching obfuscator layer, and a background. Blit the background over the obfuscator, then blit a preprocessed and alpha blended version of each of the keys over the background. This search image is saved with a unique ID and verification file of the same name containing the location of these keys is also stored.

### Display
A user is presented with the 3 normalised key images, the search image and the phrase "Find these in here". To pass the capture the user will click on the 3 locations that the keys appear within the search image.

#### Small Screen display
As the search image should be quite large small device may not have the real estate to display it and the keys together, so the display could be changed to show the first key and message "Find this in ...", then flick to the search image. When the user has made their selection, show the next key image with the text "now find this ...", then the last key and message "... and finally this."

### Verification
When the user clicks on a coordinate on the image it is stored locally. Once all 3 locations have been entered by the user they are sent to the server for verification. The locations are compared with the matched verifiaction file to see if each click is within the location of a different key.

## Why its better
Most CAPTCHA either use obfuscated text, which is solvable by modern AI faster and more effectively than humans, or a dictionary of images e.g. cats and dogs. The problem with a dictionary of images is that if the images and their categorisation are known to the attacker i.e. you want to download and install without any configuration, then it is a relatively trivial matching process to find the keys you are looking for.

FindIt gives the list of images up front, but the process of generating the search image makes it very hard for the AI and is something that the human brain has been evolved to do.

The default images of animals and jungle are not important to the process. If you have a social media site the keys could be emoticons and the backgrounds could be photos of people. If you are in the motor trade then the keys could be cars and the backgrounds empty roads. This allows the CAPTCHA to feel part of the site that it is attached to rather than something separate.

## Considerations
The alpha blending does not have to be so strong that it is hard for the human, it only needs to generate a sufficiently complex image for the AI so that mining BitCoins with the same hardware becomes more attractive.

Each key image should be distinct from all parts of the background images. E.g. having a picture of a tree as a key and then using forests as the background is going to cause problems.

When the keys are placed onto the search image they should not overlap. This makes finding the images easier for the human and verification easier for the server.