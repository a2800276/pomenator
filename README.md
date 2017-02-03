# Pomenator

This is a tiny utility I wrote for myself to deal with hellish struggle
of publishing artefacts to Maven Central. As far as I can tell, Maven
Central is the defacto standard location that Open Source stuff needs to
be published to in order for it to be usable for the community. And as
far as I can tell, Nexus / Sontype seems to be the only possibility to
access Maven Central.

I'm sure the fine folks at Sonatype are great dedicated people and are
aware of the short commings of their system, so my apologies to them for
the following. 

The system seems to be a cruel joke, it's only documented through a
series of unlinked youtube
[videos](http://central.sonatype.org/pages/producers.html) that force
you to listen to 2 minutes of awful flamenco music at the start of each
of the ~10 clips and are narrated by the slowest speaking moderator ever
who also has a very thick accent.

The process is as follows:

- create a Jira account (https://issues.sonatype.org)
- create an issue to register a new groupid (this is the groupid from
  the maven groupId - artifactId - version trifecta), this issue needs
  to be created in the Project: "Community Support - Open Source Project
  Repository Hosting"
- wait until some approves it. This seems to be automated if your email
  domain corresponds with the group id

Ok, then you need to log into:

- Log in to https://oss.sonatype.org a web frontend was obviously
  built using some sort of web gui builder that SAP discontinued in the
  early 90's.
- Don't try to figure out what any of the crap their means. Seriously.
- Click on the link "Staging Upload" (in the left navigation)
- Select the "Upload mode" you're using (more on this later, for now,
  let's assume you have a pom and a bunch of signed jars.) Select "POM"
- Now upload the POM file (this is done using semi standard html widgets
  and assigning your POM a new fake-name rooted at 'c:\'
- Now you need to upload 7 (!) further files:
  - a jar file containing your classes
  - a jar file containing your sources
  - a jar file containing the javadoc generated from the sources
  - for each of the jar files and the POM file you need to upload
    a PGP signature (*.ASC) (see below)
- You do this by selecting each of the 7 files individually and then
  clicking "Add Artifact", finally hit: "Upload Artifacts"

This is probably a good time to take up smoking. Because -- you've
guessed it -- you're not done yet.

- Click on "Staging Repositories"
- Don't get distracted by the fact that there's a million other people's
  repositories (that's just what they call them, don't try to make sense
  of the name "repository") in the list.
- Your repository, i.e. your upload is way at the bottom of the list.
  Scroll down. Keep scrolling.
- All the millions of other Repos are named "central_bundle-12345",
  yours will be call "<your_group_id>-1234"
- If you don't find your upload on the list, it's likely that the
  mainframe is down or has not processed the batch yet. You'll need to
  click on "Refresh"
- Once it shows up you can click on your project and a bunch of fancy
  html widget trees and shit appear. You can click aroung on them and
  you may find an indication of an error, in case you've missed
  something.
- If there are no errors in the fancy html tree widgets, the "Release"
  button at the top of the page should no longer be greyed out. You can
  click on that.
- Then you need to "refresh" to make sure your thingie has disappeared
  from the list.

This is probably a good time to learn assembler or take up heroin.
Because you're only almost done. You still need to go back to Jira
(https://issues.sonatype.org) and activate synchronization.  

## OMGWTFBBQ you can't be serious!?

Yes I am. Kind of. There are a number of ways to simplify this process:

### Use maven

Maven has tasks or whatever maven calls tasks to handle all of this
crap. If you're using maven already, great. No way in hell I'm
completely restructuring my project to use maven to manage my project.

### gradle or sbt or whatever

Ok. whatever

### Automate this yourself using the REST API provided by sonatype

This is what I'm trying to do here. Unfortunately -- I shit you not --
Sonatype does not provide any documentation for their REST API (unless
you have an Admin account) and their official (well blog) advice is
this: read [this blog
post](http://blog.sonatype.com/2012/07/learning-the-nexus-rest-api-read-the-docs-or-fire-up-a-browser/) which whines on about how hard it is to write
up to date documentation for a REST API and then explains how to use the
Chrome network inspector to see what calls their web frontend makes.

### User bundle uploads

In case you are stuck doing this manually, here's the single best piece
of advice I can offer you: sonatype/nexus accepts "bundle" uploads. This
is one of the alternative modes that's offered in the "Upload Mode"
chooser. This allows you to pack all 8 of the jars and signatures you
need to upload into a single jar file. Rumor has it that Sonatype
implemented this option because the cost of buying new mice was becoming
unsustainable from all the clicking involved in the process.

### use this tool

Wouldn't suggest that. It's a crappy tool I wrote for myself and your
computer will probably catch fire.

The following notes are for myself:

- create a .json file corresponding to the required pom data. A sample
  is in the `test` directory.
- additionally, the json file contains the following properties:
  - sources : an array of directories containing sources
  - classes : an array of directories containing classes
  - output  : the output and working directory where all the artifacts
    get dumped.



## GPG

Now you just need to deal with all the GPG crap. Here's the short
version:

- install GPG, you will need to figure out how to do this yourself
- create an openpgp signing key for your project. My advice is to use
  the following invovation:

    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --secret-keyring ./secring.gpg \
        --gen-key

- Note your `gpg` command may be called `gpg2`.
- `--no-default-keyring` means to not store the new key in the system
  wide keyring with all of your other gpg keys. 
- `--keyring` tells gpg where to store the public half of the key
  instead
- just as `--secret-keyring` tells gpg where to store the secret half.
  Make sure that file is added to your `.gitignore` or similar so you
  don't publish it. On the other hand, who cares? it's not like anyone
  checks signatures.
- The filenames provided must have at least one slash ('/') in them,
  else gpg will store them in some magic secret hidden directory off
  your home directory and you'll be left wondering why generating the
  key didn't work.

You'll be asked a million questions during key gneration. Most of them
are obvious. Here's my advice:  when asked about the
'kind of key' that I want:

    Please select what kind of key you want:
      (1) RSA and RSA (default)
      (2) DSA and Elgamal
      (3) DSA (sign only)
      (4) RSA (sign only)

I pick one of the "sign only" options to make sure these keys aren't
ever used to encrypt email.

When asked for the keysize, I pick the maximum available.

When asked for the validity period, I pick the default.

Now you need to figure out the key id:

    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --secret-keyring ./secring.gpg \
        --list-keys

Your keyring should only contain a single key, so you get output like:

    ./pompubring.gpg
    ----------------
    pub   3072D/4B8FD19D 2017-02-03
    uid       [ultimate] DONT TRUST THIS KEY (TEST SIGNING KEY DONT TRUST SOFTWARE SIGNED WITH THIS KEY!) <test@example.com>

This keys key id is `4B8FD19D`. You still need to publish the key:

    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --secret-keyring ./secring.gpg \
        --send-keys 4B8FD19D

Of course you need to substitute your own key id. You can also provide a
custom keyserver to publish the key on using the `--keyserver` argument
as sonatype suggests. I have no idea why they suggest this, or if they
enfore it. Anyway, just to keep you on your toes:
hkp://pool.sks-keyservers.net is the key server they'd like you to use.

Finally, to sign the jars generating the ASC files that you need for
Sonatype to accept your uploads:


    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --secret-keyring ./secring.gpg \
        -u 4B8FD19D \ 
        -ab whatever_your.jar

Obviously taking care to replace `whatever_your.jar` with the name of
the jar file you want to sign and using your own key id. This will
generate a `whatever_your.jar.asc` file.

If you'd like you can check that the generated signature is correct:

    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --verify whatever_your.jar.asc

You'll get output indicating whether the signature corresponds to the
file. It will say "BAD" in capital letters if it doesn't.





