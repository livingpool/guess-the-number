function insertIndex() {
  const bodySection = htmx.find("#popup-leaderboard tbody");
  const name = document.getElementById("form-container").getAttribute("name");

  for (let i = 0; i < bodySection.rows.length; i++) {
    const row = bodySection.rows[i];
    const newCell = row.insertCell(0);
    newCell.textContent = i + 1;

    if (name && row.cells[1].textContent.trim() == name) {
      row.style.color = "orange";
    }
  }
}

function openLeaderboardPopup() {
  document.getElementById("overlay-leaderboard").style.display = "flex";
  document.getElementById("popup-leaderboard").style.display = "flex";
}

function closeLeaderboardPopup() {
  document.getElementById("overlay-leaderboard").style.display = "none";
  document.getElementById("popup-leaderboard").style.display = "none";
}
