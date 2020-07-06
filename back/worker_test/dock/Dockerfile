FROM python
RUN pip install --upgrade pip
RUN pip install requests flask numpy opencv-python pytest
RUN pip install torch==1.5.0+cpu torchvision==0.6.0+cpu -f https://download.pytorch.org/whl/torch_stable.html
RUN mkdir test
CMD pytest ./test/test.py



