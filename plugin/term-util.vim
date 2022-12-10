if ! has("nvim")
    finish
endif
if exists('g:term_util_loaded') && g:term_util_loaded
    finish
endif

let g:term_util_loaded = 1

let s:this_folder = fnamemodify(resolve(expand('<sfile>:p')), ':h')
let s:bin_folder = printf("%s/../bin", s:this_folder)

let $PATH=printf("%s:%s", s:bin_folder, $PATH)
let $EDITOR=printf("%s/term_util_editor", s:bin_folder)
