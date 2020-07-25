<svelte:head>
	<title>Docs</title>
</svelte:head>

<div style="display:flex; flex-direction:column; padding-bottom:35px; margin-bottom:35px;border-bottom: solid #333333 1px;">
	<a href="docs/#what_is_this_site">What is this site?</a>
	<a href="docs/#languages">What languages and libraries are supported?</a>
	<a href="docs/#testing">What testing framework is used?</a>
	<a href="docs/#access_solution">How do I access the solution from the test?</a>
	<a href="docs/#limitations">What are the limitations, can my program run forever?</a>
	<a href="docs/#tech_stack">Tech stack</a>
	<a href="docs/#future">How could this site be further developped?</a>
	<a href="docs/#styling">The styling of the site is horrendous</a>
	<a href="docs/#search">Is there search functionality?</a>
	<a href="docs/#github">Why logging with Github is required to post?</a>
	<a href="docs/#third_party">Any 3rd party tracking used?</a>
	<a href="docs/#contact">Who developed this, how to contact?</a>


</div>

<a name="what_is_this_site">
	<h1>What is this site?</h1>
</a>
<p>
	Think Quora or StackOverflow except instead of posting
	questions and answers you post tests and solutions. The
	solutions are automatically run against the tests and results are reported (only final status for now- reporting on memory,
	cpu usage, latency, code size etc. are all WIP). Of course current
	features are very limited, nothing to compare with the big sites (for e.g.
	there is still no voting, commenting, reputation system etc.)
</p>

<a name="languages">
	<h1>What languages and libraries are supported?</h1>
</a>
<p>
	Currently only Python 3.8.* with a few popular libraries e.g. flask,
	numpy, opencv, pytorch (1.5.*). In the future more language could be supported
	e.g. golang, julia, rust, of course js (node) etc. and also more libraries.
</p>

<a name="testing">
	<h1>What testing framework is used?</h1>
</a>
<p>
	<a href="https://docs.pytest.org/en/latest/">the great PyTest</a>
</p>

<a name="access_solution">
	<h1>How do I access the solution from the test?</h1>
</a>
<p>
	The solution is saved as "solution.py" in the same folder as the test so you can
	"import solution as S" and then use it e.g. "assert 1 + S.one() == 1".
</p>

<a name="limitations">
	<h1>What are the limitations, can my program run forever?</h1>
</a>
<p>
	Up to 10 seconds for now. Memory, disk-size, etc. are all consrtained of course.
</p>

<a name="tech_stack">
	<h1>Tech stack</h1>
</a>
<p>
	<a href="https://svelte.dev/">Svelte</a> <a href="https://sapper.svelte.dev/">(Sapper)</a>,
	Postgres, golang (no web framework), exported static html and js hosted in Github-Pages,
	backend api, workers and DB on a $5 per month server. Code highlighting by
	<a href="https://highlightjs.org/">highlight.js</a>. Why no HTML over the wire? because it is
	easier with Svelte, it has <a href="https://svelte.dev/tutorial/nested-components">scoped css</a>, 
	it is fast, small, etc. So for dev sanity. It's really a joy to use. Why no SSR? because
	it is easier on the server just serving a bit of json without running a node process. It is
	trivial to add with Sapper, maybe when I get a bigger server.
</p>

<a name="future">
	<h1>How could this site be further developped?</h1>
</a>
<p>A few posibilities:</p>
<ul>
	<li>just for fun, enhancing with many features like commenting, voting, private tests etc.</li>
	<li>a marketplace for test and solutions: pay to see the solution for the test. Also use private  tests
	     so solutions cannot be overfitted to the tests. This would require a reputation system, as well as
		 full CI features.</li>
	<li>a tool to manage software development in a TDD style: imagine a tech lead that publishes tests
	and a team that implements solutions. This would require integration with existing CI solutions, maybe Github etc.</li>
</ul>

<a name="styling">
	<h1>The styling of the site is horrendous</h1>
</a>
<p>indeed. But functional. I'm just a programmer :)</p>

<a name="search">
	<h1>Is there search functionality?</h1>
</a>
<p>not yet. Maybe one day...</p>

<a name="github">
	<h1>Why logging with Github is required to post?</h1>
</a>
<p>to prevent abuse.</p>

<a name="third_party">
	<h1>Any 3rd party tracking used?</h1>
</a>
<p>No.</p>

<a name="contact">
	<h1>Who developed this, how to contact?</h1>
</a>
<p>Ariel Keselman (skariel@gmail.com)</p>

