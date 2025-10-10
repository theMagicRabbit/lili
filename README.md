# lili

A simple HTML parser for LinkedIn Sales Navigator lists. 

This is not a webscraper (technically), but LinkedIn likely would not agree with me.
Use of this software may violate *your* TOS with LinkedIn. I have not read the TOS nor
have I accepted them (I do not use LinkedIn. :wink: )

If LinkedIn discovers you are using this software, they could choose to ban you from
their product. Your use of LinkedIn products is between you and them. I am not responsible
for LinkedIn. Do not try to hold me responsible for what LinkedIn may do to you because
you decided to use some random software you found online.

Talk to your own legal advisors if you have concerns about this software.

tl;dr: I'm not responsible for your use of this software or LinkedIn retaliations.

(Also, I'm not affiliated with LinkedIn or Microsoft in anyway.)

## Install

1. [Install go](https://go.dev/doc/install)
2. After go is installed, open your terminal or command prompt.
3. Paste the following command `go install github.com/theMagicRabbit/lili` into your command prompt and press <ENTER>

## Running

Before you run `lili`, you need to set up the configuration file. The configuration is in TOML, which is
basically like an INI file. The instructions below are for multiple operating systems. Just skip the ones
that are not relevant to your OS.

### First time configuration (this is not as hard as it looks, it's just step by step instructions to create a file and copy content into it.)

1. Create the `lili` configuration folder at the location specified below.
    * Linux: Create a new folder `$XDG_CONFIG_HOME/lili` (If you don't know what that means, you're going to need to do some googling; you wouldn't be using Linux if you didn't want to learn, right?)
    * macOS: Create a new folder `/Users/<your_username>/Library/Application Support/lili` (I don't know much about macOS and I'm not sure what advice to give you. Open an issue if need help?)
    * Windows: Create a new folder `%AppData%\lili` (This is kind of annoying on Windows, PowerShell might be the way to go. Otherwise, past `%AppData%` into Files/File Explorer and then create a new folder in the folder that comes up.)
2. Create an empty file config file name `config.toml` in the folder you just created.
    * Make sure that the file is a `.toml` and not `.toml.txt`. Windows is particually difficult about this at times.
3. Open `config.toml` with a plain text editor. Read below if you're not sure what a plain text editor is. *Do not use a word processor.*
    * Linux: Any text editor will work, even Emacs. :wink: If you actually don't know, `GNOME Text Editor` is probably installed on your distro by default. (Unless you're on a KDE distro, in which case try KWrite. Otherwise, search your application launcher for `text editor` and use whatever comes up as long as it's not LibreOffice Writer, which is a Word Processor.)
    * macOS: I think the Apple text editor is just called `TextEdit`.
    * Windows: `Notepad`
4. On this web page, find the file called `liliconfg.sample` and open it. 
    * Pro-gamer move is to open it in a new tab.
5. Click the `Copy raw file` button. It looks like a two overlapping squares and should be to the top right of the file display box.
    * You may need to hoover your mouse cursor to see the buton name.
    * This will copy the contents of the file to your clipboard without the line numbers.
6. Paste the sample file into the file you created on your computer and have open with your text editor.
7. Customize as you see fit.
    * This means, if you'd rather store the files somewhere than my defaults you can change that here, with caveats.
    * Caveat 1: currently you have to store the file inside your home directory.
    * Caveat 2: currently paths in the configuration file will be interpeted as releative paths to your home directory, so do not put in an absolute path.
    * Yes, both of these are related (the same) issue and I'm going to fix it.
    * Caveat 3: Don't mess up the format of the file or `lili` will not work.
    * If you're not complete sure you know what you're doing, the defaults are probably fine and you can always change them later.
8. Make a note of the directories in the file. (You can always check the configuration file if you forget.)
    * `input` is where `lili` will look for `.html` files to process. 
    * `output` is where `lili` will write `.csv` files.
    * `archive` is where `lili` will move processed `.html` files.
9. Save and close the configuration file.
10. Create the `input` folder.
    * This is the same as creating any other folder on your computer.
    * Make sure the path and name is exactly the same as you put in the configuration file.
    * By default this is `Downloads/lili` becauase this is how web browsers like to store files.
    * `lili` will create the `output` and `archive` folders, so you don't need to create those yourself.
    * You have to manually create `input` since `lili` cannot process data if the data doesn't exist.
    * Ideally, the only thing in this folder will be HTML files that need to be processed by `lili`

### Using lili

0. Put HTML files into the `input` folder that you created.
1. In your terminal, type (or paste) `lili` and press <ENTER> on your keyboard

If everything was installed and configured correctly, `lili` will read html files from `input`, write `csv` files to the `output`
folder, store a backup of the html file in `archive`, and delete the original html file. You are free to delete the archive files
if you do not want/need a backup of the files you saved. (Perhaps you want to save a copy to process them again later or in case
LinkedIn changes the format of their lists and you want me to update `lili` to work with the new html files.)

## Getting help

You're welcome to open issues with requests for bug fixes or just general help if you need. If you have anything that
could point LinkedIn back to your company, I would recommend using the email feature in GitHub to reach out to me directly
just in case.

I will make best effort to keep this working, but I can't make promises beyond that.

If you find this software helpful, please feel free to star the repo, or just send me an email letting me know
that you're using it. It's always nice to know when your tools are helpful to people.

In the meantime, happy prospecting!

