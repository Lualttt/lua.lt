from flask import Flask, render_template, send_from_directory, request
from PIL import Image
import requests
import base64
import atexit
import os
import io

app = Flask(__name__)
app.jinja_env.add_extension("pypugjs.ext.jinja.PyPugJSExtension")

VISITS = "something went wrong please contact me"
try:
  with open("visits.txt", "r") as file:
    VISITS = int(file.read().rstrip())
except (ValueError, FileNotFoundError) as error:
  print("Whoopsie!!", error)
  with open("error.txt", "a+") as file:
    file.write(str(error))
  pass

with open("drawing.txt", "r") as file:
    DRAWING_ADDRESS = file.read().rstrip()

def save_visits():
  global VISITS

  with open("visits.txt", "w") as file:
    file.write(str(VISITS))

@app.route("/favicon.ico")
def favicon():
    return send_from_directory(os.path.join(app.root_path, "static/images"),
      "cherry-blossom_1f338.ico", mimetype="image/vnd.microsoft.icon")

@app.route("/")
def index():
  global VISITS

  if type(VISITS) == int:
    VISITS += 1
    if VISITS % 10 == 0:
      save_visits()

  return render_template("index.pug", visits=VISITS)

@app.route("/visits")
def visits():
    global VISITS
    return str(VISITS)

@app.route("/process", methods=["POST"])
def process():
    if "black" not in request.form or "red" not in request.form:
        return "", 400

    black_data = request.form["black"]
    red_data = request.form["red"]

    if not black_data.startswith("data:image/png;base64,") or not red_data.startswith("data:image/png;base64,"):
        return "", 400

    black_base64 = black_data.removeprefix("data:image/png;base64,")
    red_base64 = red_data.removeprefix("data:image/png;base64,")

    black_buffer = io.BytesIO(base64.b64decode(black_base64))
    red_buffer = io.BytesIO(base64.b64decode(red_base64))

    black_image = Image.open(black_buffer)
    red_image = Image.open(red_buffer)
    white_background = Image.new("RGBA", (264, 176), (255, 255, 255, 255))

    black_image = black_image.crop((0, 0, 264, 176))
    red_image = red_image.crop((0, 0, 264, 176))

    black_image = Image.alpha_composite(white_background, black_image)
    red_image = Image.alpha_composite(white_background, red_image)

    _red, green, *_ = red_image.split()
    red_image = Image.merge("RGB", (green, green, green))

    black_buffer = io.BytesIO()
    black_image.save(black_buffer, format="PNG")
    red_buffer = io.BytesIO()
    red_image.save(red_buffer, format="PNG")

    post_data = {
        "black": base64.b64encode(black_buffer.getvalue()),
        "red": base64.b64encode(red_buffer.getvalue())
    }
    requests.post(DRAWING_ADDRESS, files=post_data)

    return "", 200

atexit.register(save_visits)

if __name__ == "__main__":
  app.run(host="0.0.0.0", debug=True)
