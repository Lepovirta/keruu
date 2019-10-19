package data

const defaultCSS string = `
body {
  background-color: #f1e8fc;
  padding: 0;
  margin: 0;
  font-family: sans-serif;
}

.wrapper {
  width: 700px;
  margin: 10px auto;
}

header {
  background-color: #44118d;
  color: #eee;
  padding: 5px 20px;
  margin: 0;
}

.title {
  margin: 2px 0;
}

.description {
  margin: 2px 0;
}

.post-list {
  background-color: white;
  margin: 0;
  padding: 20px;
}

.post-list li {
  list-style: none;
  padding: 5px 0;
}

.post-title {
  display: block;
  font-size: 1.2em;
}

.post-hostname {
  font-size: 0.9em;
  color: #f44;
}

.post-time::before {
  content: '@ ';
}

.post-time {
  font-size: 0.9em;
  color: #484;
}

footer {
  border-top: 2px solid #eee;
  padding: 10px;
  font-size: 0.9em;
  text-align: right;
  color: #444;
}

@media screen and (max-width:700px) {
  .wrapper {
    width: auto;
    margin: 0;
  }
}
`
