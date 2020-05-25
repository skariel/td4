<script>
	import { onMount } from 'svelte';
	import { get, getUser } from './utils';

    let user          = {};
    let solution_id   = 0;
	let solution      = [];

	onMount(async ()=>{
        user = getUser();
		load_data()
	});

	function load_data() {
		const url = new URL(location)
        solution_id = url.searchParams.get("id")
		get(user, 'solution/'+solution_id)
			.then((r)=>{solution=r.data;})
	}
</script>

<style>
    .top {
        display:       flex;
    }
    .avatar {
        width:    40px;
        height:   40px;
        margin-right:20px;
    }

    .solutionid {
        margin-left: auto;
    }

	.code {
		background-color: #333333;
		min-height: 100px;
	}

	code {
		background-color: #00000000;
		color: antiquewhite;
	}
	.title {
		display: flex;
		align-items: center;
	}
</style>

<svelte:head>
	<title>Solution {solution_id}</title>
</svelte:head>

<div class="top">
	<img class="avatar" src={solution.avatar} alt="avatar"/>
	<h4 class="displayname">{solution.display_name}</h4>
	<h4 class="solutionid">solution {solution.id}</h4>
</div>

<h4 style="margin-top:10px;">Solution for <a href={"/test?id="+solution.test_code_id}>test {solution.test_code_id}</a></h4>

<pre class="code">
	<code>
		{solution.code}
	</code>
</pre>

<div class="title">
<!-- TODO: implement showing results of run -->
    <h1>Results: TODO!</h1>
</div>

