# Public Domain (-) 2010-2011 The Ampify Authors.
# See the Ampify UNLICENSE file for details.

# Argonought
# ==========
#
# Argonought is meant as a complementary serialisation format to JSON.

if exports?
    sys = require('sys')

divmod = (x, y) ->
    [((x - x % y) / y), x % y]

_unpack_cache = {}

unpack_number = (s) ->
    if s in _unpack_cache
        return _unpack_cache[s]
    unless s
        throw "Cannot decode an empty string"
    num = frac = 0
    ori = s
    result = ['\xff']
    num -= 8258175

if exports?
    sys.puts(divmod 63, 2)
    sys.puts(_unpack_cache['foo'])
