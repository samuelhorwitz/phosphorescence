```
with phrase as (select phraseto_tsquery('english', unaccent('change script yay foobar pure')) as phrase)
select
ts_rank(document, phrase) as rank,
searchables.id,
searchables.type,
ts_headline('english'::regconfig, name, phrase, 'StartSel = <mark>, StopSel = </mark>') as name,
ts_headline('english'::regconfig, unmodified_description, phrase, 'StartSel = <mark>, StopSel = </mark>') as description,
ts_headline('english'::regconfig, author_name, phrase, 'StartSel = <mark>, StopSel = </mark>') as author_name
from searchables, phrase
order by rank desc;
```

```
with phrase as (select phraseto_tsquery('english', unaccent('cool')) as phrase)
select rank, id, type, name, description, author_name
from (
    select ts_rank(document, phrase) as rank,
    searchables.id,
    searchables.type,
    ts_headline('english'::regconfig, name, phrase, 'StartSel = <mark>, StopSel = </mark>') as name,
    ts_headline('english'::regconfig, unmodified_description, phrase, 'StartSel = <mark>, StopSel = </mark>') as description,
    ts_headline('english'::regconfig, author_name, phrase, 'StartSel = <mark>, StopSel = </mark>') as author_name
    from searchables, phrase
) searchables
where rank > 0.01
order by rank desc;
```

```
with original_words as (select regexp_split_to_table(unaccent('cool  thß  ev ya'), '\s+') as original_words)
select array_to_string(array_agg(case when array_length(lex.corrected_words, 1) > 0 then lex.corrected_words[1] else original_words end), ' ') as fixed
from original_words
cross join lateral (
    select array(
        select word
        from searchable_lexemes
        order by word <-> original_words asc
        limit 1
    ) as corrected_words
) lex;
```

```
with original_words as (select regexp_split_to_table(unaccent('coo...  th  ever ya'), '\s+') as original_words)
select string_agg(distinct lex.word, ' ') as corrected
from original_words
cross join lateral (
    select word
    from searchable_lexemes
    where similarity(word, original_words) >= 0.3
    order by word <-> original_words asc
) lex;
```

`#((?:[\pL\pN]+(?:[\p{Pc}\p{Pd}][\pL\pN]+)*))`

```
    r := regexp.MustCompile(`#((?:[\pL\pN][\pM\x{200C}\x{200D}]*)+(?:[\p{Pc}\p{Pd}](?:[\pL\pN][\pM\x{200C}\x{200D}]*)+)*)`)
    fmt.Println(strings.Join(r.FindAllString("#füg #각 #각 #क्षि-각_üsfweg", -1), ","))
```