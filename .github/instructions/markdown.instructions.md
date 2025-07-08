---
description: 'Documentation and content creation standards'
applyTo: '**/*.md'
---

## Markdown Content Rules

The following markdown content rules are enforced in the validators:

1. **Headings**: Use appropriate heading levels (H2, H3, etc.) to structure your content. Do not use an H1 heading, as this will be generated based on the title.
2. **Lists**: Use bullet points or numbered lists for lists. Ensure proper indentation and spacing.
3. **Code Blocks**: Use fenced code blocks for code snippets. Specify the language for syntax highlighting.
4. **Links**: Use proper markdown syntax for links. Ensure that links are valid and accessible.
5. **Images**: Use proper markdown syntax for images. Include alt text for accessibility.
6. **Tables**: Use markdown tables for tabular data. Ensure proper formatting and alignment.
7. **Line Length**: Limit line length to 400 characters for readability.
8. **Whitespace**: Use appropriate whitespace to separate sections and improve readability.
9. **Front Matter**: Include YAML front matter at the beginning of the file with required metadata fields.

## Formatting and Structure

Follow these guidelines for formatting and structuring your markdown content:

- **Headings**: Use `##` for H2 and `###` for H3. Ensure that headings are used in a hierarchical manner. Recommend restructuring if content includes H4, and more strongly recommend for H5.
- **Lists**: Use `-` for bullet points and `1.` for numbered lists. Indent nested lists with two spaces.
- **Code Blocks**: Use triple backticks to create fenced code blocks. Specify the language after the opening backticks for syntax highlighting (e.g., `csharp`).
- **Links**: Use `[link text](URL)` for links. Ensure that the link text is descriptive and the URL is valid.
- **Images**: Use `![alt text](image URL)` for images. Include a brief description of the image in the alt text.
- **Tables**: Use `|` to create tables. Ensure that columns are properly aligned and headers are included.
- **Line Length**: Break lines at 80 characters to improve readability. Use soft line breaks for long paragraphs.
- **Whitespace**: Use blank lines to separate sections and improve readability. Avoid excessive whitespace.

## Follow our Guidelines

### Spelling

In cases where American spelling differs from Commonwealth/"British" spelling, use the American spelling.

Although non-American readers tend to be tolerant of reading American spelling in technical documentation,
they may find it difficult to have to type American spelling.
For example, if your documentation tells a reader who's used to the spelling colour to type color,
they may mistype it. So when you use filenames, URLs, and data parameters in examples,
try to avoid words that are spelled differently by different groups of English speakers.

### Write accessibly

#### Ease of reading

* Do not force line breaks (hard returns) within sentences and paragraphs.
  Line breaks might not work well in resized windows or with enlarged text.
* Break up walls of text to aid in scannability.
  For example, separate paragraphs, create headings, and use lists.
* Prefer short sentences.
* Define acronyms and abbreviations on first usage and if they are used infrequently.
* Place distinguishing and important information of a paragraph in the first sentence to aid in scannability.
* Use clear and direct language. Avoid the use of double negatives and exceptions in exceptions.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```markdown
A missing path will not prevent you from continuing.
```

<ul>
  <li>Double negation (missing, not)</li>
  <li>Use of future tense (will)</li>
</ul>
</td><td>

```markdown
You can continue without a path.
```

</td></tr>
</tbody></table>

#### Headings and titles

Use descriptive headings and titles because they help a reader navigate their browser and the page.
It's easier to jump between pages and sections of a page if the headings and titles are unique.

* Use a heading hierarchy.
* Do not skip levels of hierarchy (`h3` can only exist under `h2`)
* Do not use empty headings
* Use a level-1 heading for the page title.
* Use sentence casing for titles and headings.

#### Links

* Use meaningful link text. Links should make sense when read out of context.
* Do not force links to open in a new tab or window, let the reader decide how to open links.
* When possible, avoid adjacent links. Instead, put at least one character in between to separate them.
* If a link downloads a file, indicate this action and the file type in the link text.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```markdown
Use meaningful link text like described [here](https://developers.google.com/style/link-text).
Use meaningful link text. [See document.](https://developers.google.com/style/link-text)
Use meaningful link text. https://developers.google.com/style/link-text
```

</td><td>

```markdown
Use [meaningful link text](https://developers.google.com/style/link-text).
```

</td></tr>
</tbody></table>

#### Images

* When possible, use SVG images over any other format, since they are significantly lighter while having perfect information.
* For every image, provide alt text that adequately summarizes the intent of each image.
* Most of the time, do not present new information in images; always provide an equivalent text explanation with the image. There are of course exceptions for that, such as architecture diagrams, sequence diagrams etc.
* Do not repeat images.
* Avoid images of text, use text instead.

#### Tables

* Introduce tables in the text preceding the table.
* Avoid using tables to lay out pages.
* If the table contains only a single column, use a list instead.
* Do not put tables in the middle of lists or sentences.
* Sort rows in a logical order, or alphabetically if there is no logical order.

### Use the active voice

In general, use the active voice instead of the passive voice. Make it clear who is performing the action.
When using passive voice, it is easy to neglect to indicate who or what is performing the described action.
In this kind of construction, it is often hard for readers to figure out who is supposed to do something.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```markdown
The service is queried, and an acknowledgment is sent.
The service is queried by you, and an acknowledgment is sent by the server.
```

</td><td>

```markdown
Send a query to the service. The server sends an acknowledgment.
```

</td></tr>
</tbody></table>

#### Exceptions

In certain cases, it makes more sense to use the passive voice.

* To emphasize an object over an action.
* To de-emphasize a subject or actor.
* If your readers do not need to know who is responsible for the action.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```markdown
You created over 50 conflicts in the file.
```

</td><td>

```markdown
Over 50 conflicts were found in the file.
```

</td></tr>
<tr><td>

```markdown
The system saved your file.
```

</td><td>

```markdown
The file is saved.
```

</td></tr>
<tr><td>

```markdown
A system administrator purged the database in January.
```

</td><td>

```markdown
The database was purged in January.
```

</td></tr>
</tbody></table>

### Write for a global audience

* Provide context. Do not assume that the reader already knows what you're talking about.
* Avoid negative constructions when possible. Consider whether it's necessary to tell the reader what they can't do instead of what they can.
* Avoid directional language (for example, above or below) in procedural documentation.
  This increases maintenance costs and could lead to future modifications breaking the documentation.

Here are some examples.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```markdown
This document makes use of the following terms:
```

Can be substituted for a simpler verb.

</td><td>

```markdown
This document uses the following terms:
```

</td></tr>
<tr><td>

```markdown
A hybrid cloud-native DevSecOps pipeline
```

Too many nouns as modifiers of another noun. Can be broken into two parts.

</td><td>

```markdown
A cloud-native DevSecOps pipeline in a hybrid environment
```

</td></tr>
<tr><td>

```markdown
Only request one token.
```

Misplaced modifier, makes the sentence less clear and more ambiguous.

</td><td>

```markdown
Request only one token.
Request no more than one token.
Request a single token.
```

</td></tr>
<tr><td>

```markdown
If you use the term green beer in an ad, then make sure that it is targeted.
```

Here, "it is" becomes ambiguous. It could describe the green beer or the ad.

</td><td>

```markdown
If you use the term green beer in an ad, then make sure that the ad is targeted.
```

</td></tr>
</tbody></table>

#### Use present tense

In general, use present tense rather than future tense; in particular, try to avoid using _will_ where possible.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```markdown
Send a query to the service. The server will send an acknowledgment.
```

</td><td>

```markdown
Send a query to the service. The server sends an acknowledgment.
```

</td></tr>
</tbody></table>

Sometimes, of course, future tense is unavoidable because you're actually talking about the future
(for example, _This document will be outdated once PR #12345 gets merged._).
Attempting to predict the future in a document is usually a bad idea, but sometimes it's necessary.

However, the fact that the reader will be writing and running code in the future isn't a good reason to use future tense.

Also avoid the hypothetical future would—for example:

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```markdown
You can send an unsubscribe message. The server would then remove you from the mailing list.
```

</td><td>

```markdown
If you send an unsubscribe message, the server removes you from the mailing list.
```

</td></tr>
</tbody></table>

#### Use clear, precise, unambiguous language

* Use simple words. For example, do not use words like _commence_ when you mean _start_ or _begin_.
* Define abbreviations. Abbreviations can be confusing out of context, and they don't translate well.
  Spell things out whenever possible, at least the first time that you use a given term.

#### Be consistent

If you use a particular term for a particular concept in one place, then use that exact same term elsewhere, including the same capitalization.

* Use standard English word order. Sentences follow the subject + verb + object order.
* Try to keep the main subject and verb as close to the beginning of the sentence as possible.
* Use the conditional clause first. If you want to tell the audience to do something in a particular circumstance, mention the circumstance before you provide the instruction.
* Make list items consistent. Make list items parallel in structure. Be consistent in your capitalization and punctuation.
* Use consistent typographic formats. Use bold and italics consistently. Don't switch from using italics for emphasis to underlining.
* Avoid colloquialisms, idioms, or slang. Phrases like ballpark figure, back burner, or hang in there can be confusing to non-native readers.

### Describe conditions before instructions

If you want to tell the reader to do something, try to mention the circumstance, conditions, or goal before you provide the instruction.
Mentioning the circumstance first lets the reader skip the instruction if it doesn't apply.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```markdown
See [link to other document] for more information.
Click Delete if you want to delete the entire document.
Using custom domains might add noticeable latency to responses if your app is located in one of the following regions:
```

</td><td>

```markdown
For more information, see [link to other document].
To delete the entire document, click Delete.
If your app is located in one of the following regions, using custom domains might add noticeable latency to responses:
```

</td></tr>
</tbody></table>

### Use lists

Introduce a list with the appropriate context. In most cases, precede a list with an introductory sentence.

* Use simple numbered lists for steps to be performed in order.
* Nested sequential lists can detail sub-steps as well.
* Use bulleted lists when there are no sequences or options.

### Use code blocks

In most cases, precede a code sample with an introductory sentence.

* Do not use tabs to indent code; use spaces only.
* Wrap lines at 80 characters if you need to, but try to use shorter lines in code blocks.
* Specify the code block language, for syntax highlighting.
* If the code block is meant to show a command being run, prefer showing the expected output if applicable.

### Markdown guidelines

#### Add spacing to headings

Prefer spacing after `#` and newlines before and after.

```markdown
...text before.

# Heading 1

Text after...
```

#### Use lazy numbering for long lists

Markdown is smart enough to let the resulting HTML render your numbered lists correctly.
For longer lists that may change, especially long nested lists, use _lazy_ numbering.

```markdown
1.  Foo.
1.  Bar.
    1.  Barbaz.
    1.  Barbar.
1.  Baz.
```

However, if the list is small, and you don’t anticipate changing it, prefer fully numbered lists,
because it is nicer to read in source.

#### Long links

Long links make source Markdown difficult to read and break the 80 character wrapping. Wherever possible, **shorten your links**.
If it is not possible, feel free to reference links at the bottom of the paragraph instead:

```markdown
This paragraph's lines would get very long and difficult to wrap if the [full link] is included inline.

[full link]:https://www.reallylong.link/rll/BFob89Cv/Owa_TbBBi3Bn9/n5cahxQtC4TOH/afoPnUDyyOS/_8Ilq4zSBjqmo8w/j6UN1uviS9zky
```

#### Prefer lists to tables

Any tables in your Markdown should be small.
Complex, large tables are difficult to read in source and most importantly, a pain to modify later.

Lists and subheadings usually suffice to present the same information in a slightly less compact,
though much more edit-friendly way.

Here is a bad example:

```markdown
Fruit | Attribute | Notes
--- | --- | ---
Apple | [Juicy](https://example.com/SomeReallyReallyReallyReallyReallyReallyReallyReallyLongQuery), Firm, Sweet | Apples keep doctors away.
Banana | [Convenient](https://example.com/SomeDifferentReallyReallyReallyReallyReallyReallyReallyReallyLongQuery), Soft, Sweet | Contrary to popular belief, most apes prefer mangoes.
```

And here is a better alternative:

```markdown
## Fruits

### Apple

* [Juicy](https://SomeReallyReallyReallyReallyReallyReallyReallyReallyReallyReallyReallyReallyReallyReallyReallyReallyLongURL)
* Firm
* Sweet

Apples keep doctors away.

### Banana

* [Convenient](https://example.com/SomeDifferentReallyReallyReallyReallyReallyReallyReallyReallyLongQuery)
* Soft
* Sweet

Contrary to popular belief, most apes prefer mangoes.
```

#### Strongly prefer Markdown to HTML

Please prefer standard Markdown syntax wherever possible and avoid HTML hacks.
If you can not seem to accomplish what you want, reconsider whether you really need it.
Except for big tables, Markdown meets almost all needs already.

Every bit of HTML or Javascript hacking reduces the readability and portability.
This in turn limits the usefulness of integrations with other tools, which may either present the source as plain text or render it.

#### Spacing

* Remove all trailing whitespaces at end of lines.
* Remove instances of multiple consecutive blank lines.
* Files should end with a single newline character.


## Validation Requirements

Ensure compliance with the following validation requirements:

- **Front Matter**: Include the following fields in the YAML front matter:

    - `post_title`: The title of the post.
    - `author1`: The primary author of the post.
    - `post_slug`: The URL slug for the post.
    - `microsoft_alias`: The Microsoft alias of the author.
    - `featured_image`: The URL of the featured image.
    - `categories`: The categories for the post. These categories must be from the list in /categories.txt.
    - `tags`: The tags for the post.
    - `ai_note`: Indicate if AI was used in the creation of the post.
    - `summary`: A brief summary of the post. Recommend a summary based on the content when possible.
    - `post_date`: The publication date of the post.

- **Content Rules**: Ensure that the content follows the markdown content rules specified above.
- **Formatting**: Ensure that the content is properly formatted and structured according to the guidelines.
- **Validation**: Run the validation tools to check for compliance with the rules and guidelines.

## Admonitions

Use GitHub-flavored markdown for admonitions: NOTE, WARNING, TIP, IMPORTANT, CAUTION.

Examples:

```markdown
> [!NOTE]  
> Highlights information that users should take into account, even when skimming.

> [!TIP]
> Optional information to help a user be more successful.

> [!IMPORTANT]  
> Crucial information necessary for users to succeed.

> [!WARNING]  
> Critical content demanding immediate user attention due to potential risks.

> [!CAUTION]
> Negative potential consequences of an action.
```
