<svelte:head>
	<title>Docs</title>
</svelte:head>

<h1>What is this site?</h1>
<p>
	Think Quora or StackOverflow except instead of posting
	questions and answers you post tests and solutions. The
	solutions are automatically run against the tests and results are reported (only final status for now- reporting on memory,
	cpu usage, latency, code size etc. are all WIP). Of course current
	features are very limited, nothing to compare with the big sites (for e.g.
	there is still no voting, commenting, reputation system etc.)
</p>

<h1>What languages and libraries are supported?</h1>
<p>
	Currently only Python 3.8.* with a few popular libraries e.g. flask,
	numpy, opencv, pytorch (1.5.*). In the future more language could be supported
	e.g. golang, julia, rust, of course js (node) etc. and also more libraries.
</p>

<h1>What testing framework is used?</h1>
<p>
	<a href="https://docs.pytest.org/en/latest/">the great PyTest</a>
</p>

<h1>How do I access the solution from the test?</h1>
<p>
	The solution is saved as "solution.py" in the same folder as the test so you can
	"import solution as S" and then use it e.g. "assert 1 + S.one() == 1".
</p>

<h1>What are the limitations, can my program run forever?</h1>
<p>
	Up to 10 seconds for now. Memory, disk-size, etc. are all consrtained of course.
</p>

<h1>Tech stack</h1>
<p>
	<a href="https://svelte.dev/">Svelte</a> <a href="https://sapper.svelte.dev/">(Sapper)</a>, Postgres, golang (no web framework), exported static html and js hosted in Github-Pages,
	backend api, workers and DB on a $5 per month server. Why no HTML over the wire? because it is
	easier with Svelte, it has <a href="https://svelte.dev/tutorial/nested-components">scoped css</a>, 
	it is fast, small, etc. So for the dev commodity. It's really a joy to use.
</p>

<h1>How could this site be further developped?</h1>
<p>A few posibilities:</p>
<ul>
	<li>just for fun, enhancing with many features like commenting, voting, private tests etc.</li>
	<li>a marketplace for test and solutions: pay to see the solution for the test. Also use private  tests
	     so solutions cannot be overfitted to the tests. This would require a reputation system, as well as
		 full CI features.</li>
	<li>a tool to manage software development in a TDD style: imagine a tech lead that publishes tests
	and a team that implements solutions. This would require integration with existing CI solutions, maybe Github etc.</li>
</ul>

<h1>Who developed this, how to contact?</h1>
<p>Ariel Keselman (skariel@gmail.com)</p>

