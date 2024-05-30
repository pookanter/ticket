import { Tabs as TabsPrimitive } from 'bits-ui';
import Content from './board-tabs-content.svelte';
import List from './board-tabs-list.svelte';
import Trigger from './board-tabs-trigger.svelte';

const Root = TabsPrimitive.Root;

export {
	Root,
	Content,
	List,
	Trigger,
	//
	Root as Tabs,
	Content as TabsContent,
	List as TabsList,
	Trigger as TabsTrigger
};
