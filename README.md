# Go Simple CRUD

**DON'T USE THIS, WIP.**

## Dev Notes

- Escape everything that might come out as an error from the database.
- Create tests to cover all of those instances, as yesterday it's not
  known that `json.Unmarshall` cannot accept `\\n` or `\\r` and doesn't
  automatically escape characters such as double quotes.
