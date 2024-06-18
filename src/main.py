from flask import Flask, render_template, url_for
import atexit

app = Flask(__name__)
app.jinja_env.add_extension("pypugjs.ext.jinja.PyPugJSExtension")

VISITS = "something went wrong please contact me"
try:
  with open("visits.txt", "r") as file:
    VISITS = int(file.read())
except (ValueError, FileNotFoundError) as error:
  print("Whoopsie!!", error)
  with open("error.txt", "a+") as file:
    file.write(str(error))
  pass

def save_visits():
  global VISITS

  with open("visits.txt", "w") as file:
    file.write(str(VISITS))

@app.route("/")
def index():
  global VISITS

  if type(VISITS) == int:
    VISITS += 1
    if VISITS % 10 == 0:
      save_visits()

  return render_template("index.pug", visits=VISITS)

atexit.register(save_visits)

if __name__ == "__main__":
  app.run(debug=True)
