% -*- mode: prolog -*-

member_(_, A, A).
member_([C|A], B, _) :- member_(A, B, C).

member(B, [C|A]) :- member_(A, B, C).

prefix([], _).
prefix([X|Xs], [X|Ys]) :- prefix(Xs, Ys).
    
