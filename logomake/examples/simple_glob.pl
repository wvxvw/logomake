% -*- mode: prolog -*-

not_tests(Sources) :-
    findall(CFile, (
                glob("./*.c", CFiles),
                member(CFile, CFiles),
                \+ prefix(CFile, "test_")
            ),
            Sources).

all :-
    not_tests(Sources),
    printf('Building %s~n', Sources),
    c(Sources, "program").
