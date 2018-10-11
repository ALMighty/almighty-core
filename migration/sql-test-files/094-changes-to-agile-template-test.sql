-- create space template
INSERT INTO space_templates (id,name,description) VALUES('f06ba0ba-eaf0-4655-bfe8-7b9e26d0f48f', 'test space template', 'test template');

-- create space
insert into spaces (id,name,space_template_id) values ('fe8e7e07-a8d7-41c2-9761-ca1ffe2409b4', 'test space', 'f06ba0ba-eaf0-4655-bfe8-7b9e26d0f48f');

-- create work item type Theme
INSERT INTO work_item_types (id,name,space_template_id,fields,description,icon) VALUES ('5182fc8c-b1d6-4c3d-83ca-6a3c781fa18a', 'Theme', 'f06ba0ba-eaf0-4655-bfe8-7b9e26d0f48f', '{"system_area": {"Type": {"Kind": "area"}, "Label": "Area", "Required": false, "Description": "The area to which the work item belongs"}, "system_order": {"Type": {"Kind": "float"}, "Label": "Execution Order", "Required": false, "Description": "Execution Order of the workitem."}, "system_state": {"Type": {"Kind": "enum", "Values": ["new", "open", "in progress", "resolved", "closed"], "BaseType": {"Kind": "string"}}, "Label": "State", "Required": true, "Description": "The state of the work item"}, "system_title": {"Type": {"Kind": "string"}, "Label": "Title", "Required": true, "Description": "The title text of the work item"}, "system_labels": {"Type": {"Kind": "list", "ComponentType": {"Kind": "label"}}, "Label": "Labels", "Required": false, "Description": "List of labels attached to the work item"}, "business_value": {"Type": {"Kind": "integer"}, "Label": "Business Value", "Required": false, "Description": "The business value of this work item."}, "effort": {"Type": {"Kind": "float"}, "Label": "Effort", "Required": false, "Description": "The effort that was given to this workitem within its space."}, "time_criticality": {"Type": {"Kind": "float"}, "Label": "Time Criticality", "Required": false, "Description": "The time criticality that was given to this workitem within its space."}, "system_creator": {"Type": {"Kind": "user"}, "Label": "Creator", "Required": true, "Description": "The user that created the work item"}, "system_codebase": {"Type": {"Kind": "codebase"}, "Label": "Codebase", "Required": false, "Description": "Contains codebase attributes to which this WI belongs to"}, "system_assignees": {"Type": {"Kind": "list", "ComponentType": {"Kind": "user"}}, "Label": "Assignees", "Required": false, "Description": "The users that are assigned to the work item"}, "system_iteration": {"Type": {"Kind": "iteration"}, "Label": "Iteration", "Required": false, "Description": "The iteration to which the work item belongs"}, "system_created_at": {"Type": {"Kind": "instant"}, "Label": "Created at", "Required": false, "Description": "The date and time when the work item was created"}, "system_updated_at": {"Type": {"Kind": "instant"}, "Label": "Updated at", "Required": false, "Description": "The date and time when the work item was last updated"}, "system_description": {"Type": {"Kind": "markup"}, "Label": "Description", "Required": false, "Description": "A descriptive text of the work item"}, "system_remote_item_id": {"Type": {"Kind": "string"}, "Label": "Remote item", "Required": false, "Description": "The ID of the remote work item"}}', 'Description for Planner Item', 'fa fa-bookmark');

-- create work item type Epic
INSERT INTO work_item_types (id,name,space_template_id,fields,description,icon) VALUES ('2c169431-a55d-49eb-af74-cc19e895356f', 'Epic', 'f06ba0ba-eaf0-4655-bfe8-7b9e26d0f48f', '{"system_area": {"Type": {"Kind": "area"}, "Label": "Area", "Required": false, "Description": "The area to which the work item belongs"}, "system_order": {"Type": {"Kind": "float"}, "Label": "Execution Order", "Required": false, "Description": "Execution Order of the workitem."}, "system_state": {"Type": {"Kind": "enum", "Values": ["new", "open", "in progress", "resolved", "closed"], "BaseType": {"Kind": "string"}}, "Label": "State", "Required": true, "Description": "The state of the work item"}, "system_title": {"Type": {"Kind": "string"}, "Label": "Title", "Required": true, "Description": "The title text of the work item"}, "system_labels": {"Type": {"Kind": "list", "ComponentType": {"Kind": "label"}}, "Label": "Labels", "Required": false, "Description": "List of labels attached to the work item"}, "component": {"Type": {"Kind": "string"}, "Label": "Component", "Required": false, "Description": "The component value of this work item."}, "business_value": {"Type": {"Kind": "integer"}, "Label": "Business Value", "Required": false, "Description": "The business value of this work item."}, "effort": {"Type": {"Kind": "float"}, "Label": "Effort", "Required": false, "Description": "The effort that was given to this workitem within its space."}, "time_criticality": {"Type": {"Kind": "float"}, "Label": "Time Criticality", "Required": false, "Description": "The time criticality that was given to this workitem within its space."}, "system_creator": {"Type": {"Kind": "user"}, "Label": "Creator", "Required": true, "Description": "The user that created the work item"}, "system_codebase": {"Type": {"Kind": "codebase"}, "Label": "Codebase", "Required": false, "Description": "Contains codebase attributes to which this WI belongs to"}, "system_assignees": {"Type": {"Kind": "list", "ComponentType": {"Kind": "user"}}, "Label": "Assignees", "Required": false, "Description": "The users that are assigned to the work item"}, "system_iteration": {"Type": {"Kind": "iteration"}, "Label": "Iteration", "Required": false, "Description": "The iteration to which the work item belongs"}, "system_created_at": {"Type": {"Kind": "instant"}, "Label": "Created at", "Required": false, "Description": "The date and time when the work item was created"}, "system_updated_at": {"Type": {"Kind": "instant"}, "Label": "Updated at", "Required": false, "Description": "The date and time when the work item was last updated"}, "system_description": {"Type": {"Kind": "markup"}, "Label": "Description", "Required": false, "Description": "A descriptive text of the work item"}, "system_remote_item_id": {"Type": {"Kind": "string"}, "Label": "Remote item", "Required": false, "Description": "The ID of the remote work item"}}', 'Description for Planner Item', 'fa fa-bookmark');

-- create work item type Story
INSERT INTO work_item_types (id,name,space_template_id,fields,description,icon) VALUES ('6ff83406-caa7-47a9-9200-4ca796be11bb', 'Story', 'f06ba0ba-eaf0-4655-bfe8-7b9e26d0f48f', '{"system_area": {"Type": {"Kind": "area"}, "Label": "Area", "Required": false, "Description": "The area to which the work item belongs"}, "system_order": {"Type": {"Kind": "float"}, "Label": "Execution Order", "Required": false, "Description": "Execution Order of the workitem."}, "system_state": {"Type": {"Kind": "enum", "Values": ["new", "open", "in progress", "resolved", "closed"], "BaseType": {"Kind": "string"}}, "Label": "State", "Required": true, "Description": "The state of the work item"}, "system_title": {"Type": {"Kind": "string"}, "Label": "Title", "Required": true, "Description": "The title text of the work item"}, "system_labels": {"Type": {"Kind": "list", "ComponentType": {"Kind": "label"}}, "Label": "Labels", "Required": false, "Description": "List of labels attached to the work item"}, "effort": {"Type": {"Kind": "float"}, "Label": "Effort", "Required": false, "Description": "The effort that was given to this workitem within its space."}, "system_creator": {"Type": {"Kind": "user"}, "Label": "Creator", "Required": true, "Description": "The user that created the work item"}, "system_codebase": {"Type": {"Kind": "codebase"}, "Label": "Codebase", "Required": false, "Description": "Contains codebase attributes to which this WI belongs to"}, "system_assignees": {"Type": {"Kind": "list", "ComponentType": {"Kind": "user"}}, "Label": "Assignees", "Required": false, "Description": "The users that are assigned to the work item"}, "system_iteration": {"Type": {"Kind": "iteration"}, "Label": "Iteration", "Required": false, "Description": "The iteration to which the work item belongs"}, "system_created_at": {"Type": {"Kind": "instant"}, "Label": "Created at", "Required": false, "Description": "The date and time when the work item was created"}, "system_updated_at": {"Type": {"Kind": "instant"}, "Label": "Updated at", "Required": false, "Description": "The date and time when the work item was last updated"}, "system_description": {"Type": {"Kind": "markup"}, "Label": "Description", "Required": false, "Description": "A descriptive text of the work item"}, "system_remote_item_id": {"Type": {"Kind": "string"}, "Label": "Remote item", "Required": false, "Description": "The ID of the remote work item"}}', 'Description for Planner Item', 'fa fa-bookmark');

-- create a work items for Theme (removed fields 'effort' [float], 'business_value' [integer], 'time_criticality' [float])
insert into work_items (id, type, space_id, fields) values ('cf84c888-ac28-493d-a0cd-978b78568040', '5182fc8c-b1d6-4c3d-83ca-6a3c781fa18a', 'fe8e7e07-a8d7-41c2-9761-ca1ffe2409b4', '{"system_title":"Work item 1", "effort":12.34, "business_value":1234, "time_criticality":56.78}'::json);
insert into work_items (id, type, space_id, fields) values ('8bbb542c-4f5c-44bb-9272-e1a8f24e6eb2', '5182fc8c-b1d6-4c3d-83ca-6a3c781fa18a', 'fe8e7e07-a8d7-41c2-9761-ca1ffe2409b4', '{"system_title":"Work item 2"}'::json);

-- create a work items for Epic (removed fields 'effort' [float], 'business_value' [integer], 'time_criticality' [float], 'component' [string])
insert into work_items (id, type, space_id, fields) values ('4aebb314-a8c1-4e9c-96b6-074769d16934', '2c169431-a55d-49eb-af74-cc19e895356f', 'fe8e7e07-a8d7-41c2-9761-ca1ffe2409b4', '{"system_title":"Work item 3", "effort":12.34, "business_value":1234, "time_criticality":56.78, "component":"Component 1"}'::json);
insert into work_items (id, type, space_id, fields) values ('9c53fb2b-c6af-48a1-bef1-6fa547ea72fa', '2c169431-a55d-49eb-af74-cc19e895356f', 'fe8e7e07-a8d7-41c2-9761-ca1ffe2409b4', '{"system_title":"Work item 4"}'::json);

-- create a work items for Story (removed fields 'effort' [float])
insert into work_items (id, type, space_id, fields) values ('68f83154-8d76-49c1-8be0-063ce90f803d', '6ff83406-caa7-47a9-9200-4ca796be11bb', 'fe8e7e07-a8d7-41c2-9761-ca1ffe2409b4', '{"system_title":"Work item 5", "effort":12.34}'::json);
insert into work_items (id, type, space_id, fields) values ('17e2081f-812d-4f4e-9c51-c537406bd1d8', '6ff83406-caa7-47a9-9200-4ca796be11bb', 'fe8e7e07-a8d7-41c2-9761-ca1ffe2409b4', '{"system_title":"Work item 6"}'::json);
