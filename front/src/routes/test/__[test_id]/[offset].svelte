<script context="module">
	export function preload({ params }) {
        return {
			test_id: params.test_id,
			offset  : params.offset,
		}
	}
</script>

<script>
	import { onMount, getContext } from 'svelte';
	import { get } from '../../utils';

    export let test_id;
	export let offset;
	
	const user  = getContext('user');
	let test  = [];
	let solutions = [];

	onMount(initial_load);

	function initial_load() {
		get(user, 'test/'+test_id)
			.then((r)=>{test=r.data;})
		get(user, 'solutions_by_test/'+test_id+'/'+offset)
			.then((r)=>{solutions=r.data;})
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

    .displayname {

    }

    .testid {
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

</style>

<title>test {test_id}</title>

<div class="top">
	<img class="avatar" src={test.avatar} alt="avatar"/>
	<h4 class="displayname">{test.display_name}</h4>
	<h4 class="testid">{test.id}</h4>
</div>

<pre class="code">
	<code>
		{test.code}
	</code>
</pre>

<!-- <h4 class="updated">updated: {test.ts_updated}</h4> -->

<a href={'/test/'+test_id+'/new_solution'}>add solution</a>


<h4>solutions:</h4>

{#each solutions as s }
	<div  style="margin:15px;">
		<h5>s:</h5>
		<pre>
			<code>
				{s.code}
			</code>
		</pre>
		<a href="/test/{test_id}/code/{s.id}">code id: {s.id}</a>
	</div>
{/each}
