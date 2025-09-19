const lightIcon = document.getElementById("sun");
const darkIcon = document.getElementById("moon");

// check if dark mode is preferred in user's OS
const darkModeMediaQuery = window.matchMedia("(prefer-color-scheme: dark)");
let darkMode = darkModeMediaQuery.matches;
let savedTheme = darkMode ? "dark" : "light";

// <html>.dataset is used to save the theme
const htmlEl = document.documentElement;

// toggle the theme and update their local storage.
function toggleTheme(bool) {
  const theme = bool ? "light" : "dark";
  htmlEl.dataset.theme = theme;
  localStorage.setItem("savedTheme", theme);
}

// handle their saved preferred theme.
function setSavedTheme() {
  // If there is no current theme in localstorage then give 'em the OS setting as default.
  // Like for first time visitors...
  if (localStorage.getItem("savedTheme") === null) {
    localStorage.setItem("savedTheme", savedTheme);
  } else {
    savedTheme = localStorage.getItem("savedTheme");
  }

  htmlEl.dataset.theme = savedTheme;
}

// set the default theme.
setSavedTheme();
