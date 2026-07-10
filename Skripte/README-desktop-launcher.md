Desktop-Starter fuer Asasel Quick

Datei im Repo:
- Skripte/asasel-quick.desktop

So verwenden:
1. In Skripte/asasel-quick.desktop die Platzhalter anpassen:
   - DEINUSER -> dein Linux-Account fuer die Restzeit-Abfrage
   - REPO_PATH -> absoluter Pfad zu deinem lokalen Repo, z. B. /home/ber/Projekte/Asasel
2. Datei kopieren nach:
   ~/.local/share/applications/asasel-quick.desktop
3. Ausfuehrbar machen:
   chmod +x ~/.local/share/applications/asasel-quick.desktop
4. Starter testen:
   gtk-launch asasel-quick

Hinweis:
- Als Icon wird frontend/static/favicon.png aus dem Repo verwendet.
- Falls du das Repo verschiebst, den Icon-Pfad in der .desktop-Datei entsprechend aktualisieren.
