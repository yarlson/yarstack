export function text(value, { bold = false, em = false, link = null } = {}) {
  const node = { type: 'text', text: value };
  const marks = [];
  if (bold) marks.push({ type: 'strong' });
  if (em) marks.push({ type: 'em' });
  if (link) marks.push({ type: 'link', attrs: { href: link } });
  if (marks.length) node.marks = marks;
  return node;
}

export function para(...nodes) {
  return { type: 'paragraph', content: nodes };
}

export function heading(value, level = 2) {
  return { type: 'heading', attrs: { level }, content: [text(value)] };
}

export function listItem(...nodes) {
  return { type: 'listItem', content: [para(...nodes)] };
}

export function orderedList(items) {
  return { type: 'orderedList', attrs: { order: 1 }, content: items };
}

export function bulletList(items) {
  return { type: 'bulletList', content: items };
}

export function acItem(lead, rest) {
  return listItem(text(`${lead}: `, { bold: true }), text(rest));
}

export function doc(contextNodes, acceptanceCriteriaItems, engineeringNotesNodes) {
  return {
    type: 'doc',
    version: 1,
    content: [
      heading('Context'),
      ...contextNodes,
      heading('Acceptance criteria'),
      orderedList(acceptanceCriteriaItems),
      heading('Engineering notes'),
      ...engineeringNotesNodes,
    ],
  };
}
