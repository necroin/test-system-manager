window.request_url = "%s"

function request(methood, url, data) {
    var req = new XMLHttpRequest();
    if (!url.startsWith("http://") && !url.startsWith("http//")) {
        url = "http://" + url
    }
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function post_request(url, data) {
    return request("POST", url, data)
}

function OpenPage(path) {
    window.location.href = "http://" + window.request_url + path
}

function GetSearchLineText() {
    let searchLine = document.getElementById("search-line-text")
    return searchLine.value
}

function open_dialog(dialog, overlay) {
    document.getElementById(dialog).style.display = "flex";
    document.getElementById(overlay).style.display = "block";
}

function close_dialog(dialog, overlay) {
    document.getElementById(dialog).style.display = "none";
    document.getElementById(overlay).style.display = "none";
}

function AddProject() {
    let name = document.getElementById("create-dialog-name").value
    if (name != "") {
        let response = post_request(window.location.href + "/insert", name)
        GetProjects()
        close_dialog('create-dialog', 'create-dialog-overlay')
    }
}

function AddTestCase(projectId) {
    let name = document.getElementById("create-dialog-name").value
    if (name != "") {
        let response = post_request(window.location.href + "/insert", name)
        GetTestCases(projectId)
        close_dialog('create-dialog', 'create-dialog-overlay')
    }
}

function AddTestPlan(projectId) {
    let name = document.getElementById("create-dialog-name").value
    if (name != "") {
        let response = post_request(window.location.href + "/insert", name)
        GetTestPlans(projectId)
        close_dialog('create-dialog', 'create-dialog-overlay')
    }
}

function AddProjectTag(projectId) {
    let tag = document.getElementById("settings-tags-input").value
    if (tag != "") {
        let response = post_request(window.request_url + "/project/" + projectId + "/tags/insert", tag)
        GetProjectSettings(projectId)
    }
}

function DeleteProjectTag(projectId, tag) {
    let response = post_request(window.request_url + "/project/" + projectId + "/tags/delete", tag)
    GetProjectSettings(projectId)
}

function GetProjectTags(projectId) {
    let response = post_request(window.request_url + "/project/" + projectId + "/tags/get")
    let data = JSON.parse(response)
    return data.records
}

function GetProjects() {
    let projectsList = document.getElementById("projects");
    projectsList.replaceChildren();
    let response = post_request(
        window.request_url + "/projects/get",
        JSON.stringify({ "search": GetSearchLineText() })
    )
    let records = JSON.parse(response).records;
    for (recordIndex in records) {
        let record = records[recordIndex]
        let element = document.createElement("div");
        element.className = "list-item"
        let id = document.createElement("span");
        let name = document.createElement("span");
        let count = document.createElement("span");
        id.innerText = record.fields.Id;
        name.innerText = record.fields.Name;
        count.innerText = record.fields.TestCaseCount;
        count.style.minWidth = 100;
        count.style.maxWidth = 100;

        element.appendChild(id);
        element.appendChild(name);
        element.appendChild(count);

        let projectId = id.innerText
        element.onclick = () => OpenPage("/project/" + projectId + "/cases");

        let tagsElement = document.createElement("div")

        let tags = GetProjectTags(projectId)
        for (tagIndex in tags) {
            let tag = tags[tagIndex]
            let tagElement = document.createElement("span")
            tagElement.innerText = tag.fields.Name
            tagsElement.appendChild(tagElement)
        }
        element.appendChild(tagsElement);

        projectsList.appendChild(element);
    }
}

function RenameProject(projectId) {
    let name = document.getElementById("settings-name-input").value
    let response = post_request(window.request_url + "/project/" + projectId + "/rename", name)
    location.reload()
}

function AddProjectUser(projectId) {
    let user = document.getElementById("settings-collaborators-input").value
    let response = post_request(window.request_url + "/project/" + projectId + "/collaborators/add", user)
    location.reload()
}

function DeleteProjectUser(projectId) {
    let user = document.getElementById("settings-collaborators-input").value
    let response = post_request(window.request_url + "/project/" + projectId + "/collaborators/delete", user)
    location.reload()
}

function AddTestCaseTag() {
    let tag = document.getElementById("settings-tags-input").value
    if (tag != "") {
        let response = post_request(window.location + "/tags/insert", tag)
        UpdateTestCaseTags()
    }
}

function DeleteTestCaseTag(tag) {
    let response = post_request(window.location + "/tags/delete", tag)
    UpdateTestCaseTags()
}

function GetTestCaseTags() {
    let response = post_request(window.location + "/tags/get")
    let data = JSON.parse(response)
    return data.records
}


function AddTestPlanTag() {
    let tag = document.getElementById("settings-tags-input").value
    if (tag != "") {
        let response = post_request(window.location + "/tags/insert", tag)
        UpdateTestPlanTags()
    }
}

function DeleteTestPlanTag(tag) {
    let response = post_request(window.location + "/tags/delete", tag)
    UpdateTestPlanTags()
}

function GetTestPlanTags() {
    let response = post_request(window.location + "/tags/get")
    let data = JSON.parse(response)
    return data.records
}

function GetTestCases(projectId) {
    let testCasesList = document.getElementById("test-cases");
    testCasesList.replaceChildren();
    let response = post_request(window.request_url + "/project/" + projectId + "/cases/get",
        JSON.stringify({ "search": GetSearchLineText() }));
    let records = JSON.parse(response).records;
    for (recordIndex in records) {
        let record = records[recordIndex]
        console.log(record)
        let element = document.createElement("div");
        element.className = "list-item"
        let id = document.createElement("span");
        let name = document.createElement("span");
        id.innerText = record.fields.Id;
        name.innerText = record.fields.Name;

        element.appendChild(id);
        element.appendChild(name);

        let testCaseId = id.innerHTML
        element.onclick = () => OpenPage("/project/" + projectId + "/case/" + testCaseId);

        var tagsElement = document.createElement("div")

        let tag_response = post_request(window.request_url + "/project/" + projectId + "/case/" + testCaseId + "/tags/get")
        let tags = JSON.parse(tag_response).records

        for (tagIndex in tags) {
            let tag = tags[tagIndex]
            let tagElement = document.createElement("span")
            tagElement.innerText = tag.fields.Name
            tagsElement.appendChild(tagElement)
        }
        element.appendChild(tagsElement);

        testCasesList.appendChild(element);
    }
}

function GetTestPlans(projectId) {
    let testCasesList = document.getElementById("test-plans");
    testCasesList.replaceChildren();
    let response = post_request(window.request_url + "/project/" + projectId + "/plans/get",
        JSON.stringify({ "search": GetSearchLineText() }));
    let records = JSON.parse(response).records;
    for (recordIndex in records) {
        let record = records[recordIndex]
        let element = document.createElement("div");
        element.className = "list-item"
        let id = document.createElement("span");
        let name = document.createElement("span");
        let count = document.createElement("span");
        id.innerHTML = record.fields.Id;
        name.innerHTML = record.fields.Name;
        count.innerHTML = record.fields.TestCaseCount;

        element.appendChild(id);
        element.appendChild(name);
        element.appendChild(count);

        let testPlanId = id.innerHTML
        element.onclick = () => OpenPage("/project/" + projectId + "/plan/" + testPlanId);


        var tagsElement = document.createElement("div")

        let tag_response = post_request(window.request_url + "/project/" + projectId + "/plan/" + testPlanId + "/tags/get")
        let tags = JSON.parse(tag_response).records

        for (tagIndex in tags) {
            let tag = tags[tagIndex]
            let tagElement = document.createElement("span")
            tagElement.innerText = tag.fields.Name
            tagsElement.appendChild(tagElement)
        }
        element.appendChild(tagsElement);

        testCasesList.appendChild(element);
    }
}

function GetProjectCollaborators(projectId) {
    let response = post_request(window.request_url + "/project/" + projectId + "/collaborators")
    return JSON.parse(response)
}

function GetProjectSettings(projectId) {
    let tagsList = document.getElementById("tags");
    tagsList.replaceChildren()
    let tags = GetProjectTags(projectId)
    for (tagIndex in tags) {
        let tag = tags[tagIndex]
        let div = document.createElement("div")

        let tagElement = document.createElement("span")
        tagElement.innerText = tag.fields.Name

        let deleteButton = document.createElement("button")
        deleteButton.innerText = "✖"
        deleteButton.onclick = () => { DeleteProjectTag(projectId, tagElement.innerText) }

        div.appendChild(tagElement)
        div.appendChild(deleteButton)
        tagsList.appendChild(div)
    }

    let collaboratorsList = document.getElementById("collaborators")
    collaboratorsList.replaceChildren()
    let collaborators = GetProjectCollaborators(projectId).records
    for (index in collaborators) {
        let collaborator = collaborators[index]

        if (collaborator.fields.Role != "Создатель") {
            let listElement = document.createElement("div")
            let nameElement = document.createElement("span")
            nameElement.innerText = collaborator.fields.Username

            let roleSelect = document.createElement("select")

            let guestOption = document.createElement("option")
            let analystOption = document.createElement("option")
            let testerOption = document.createElement("option")

            guestOption.innerText = "Гость"
            analystOption.innerText = "Аналитик"
            testerOption.innerText = "Тестировщик"

            roleSelect.appendChild(guestOption)
            roleSelect.appendChild(analystOption)
            roleSelect.appendChild(testerOption)

            listElement.appendChild(nameElement)
            listElement.appendChild(roleSelect)


            if (collaborator.fields.Role == "Аналитик") {
                roleSelect.selectedIndex = 1
            }

            if (collaborator.fields.Role == "Тестировщик") {
                roleSelect.selectedIndex = 2
            }

            let deleteButton = document.createElement("button")
            deleteButton.innerText = "✖"
            listElement.appendChild(deleteButton)

            collaboratorsList.appendChild(listElement)
        }
    }
}

function UpdateTestCaseTags() {
    let tagsList = document.getElementById("tags");
    tagsList.replaceChildren()
    let tags = GetTestCaseTags()

    for (tagIndex in tags) {
        let tag = tags[tagIndex]
        let div = document.createElement("div")

        let tagElement = document.createElement("span")
        tagElement.innerText = tag.fields.Name

        let deleteButton = document.createElement("button")
        deleteButton.innerText = "✖"
        deleteButton.onclick = () => { DeleteTestCaseTag(tagElement.innerText) }

        div.appendChild(tagElement)
        div.appendChild(deleteButton)
        tagsList.appendChild(div)
    }
}

function UpdateTestPlanTags() {
    let tagsList = document.getElementById("tags");
    tagsList.replaceChildren()
    let tags = GetTestPlanTags()

    for (tagIndex in tags) {
        let tag = tags[tagIndex]
        let div = document.createElement("div")

        let tagElement = document.createElement("span")
        tagElement.innerText = tag.fields.Name

        let deleteButton = document.createElement("button")
        deleteButton.innerText = "✖"
        deleteButton.onclick = () => { DeleteTestPlanTag(tagElement.innerText) }

        div.appendChild(tagElement)
        div.appendChild(deleteButton)
        tagsList.appendChild(div)
    }
}

function GetTestCase() {
    let response = post_request(window.location.href + "/get");

    let data = JSON.parse(response)
    let description = document.getElementById("description-text")
    let scenario = document.getElementById("scenario-text")

    if (data.description != null) {
        description.innerText = data.description
    }

    if (data.scenario != null) {
        scenario.innerText = data.scenario
    }

    UpdateTestCaseTags()

}

function UpdateTestPlatTestCasesList(testCasesList) {
    let data = { "cases": [] }
    for (i = 0; i < testCasesList.children.length; i++) {
        let testCase = testCasesList.children.item(i)
        let id = testCase.children.item(0).innerText
        data["cases"].push({ "id": id })
    }
    return post_request(window.location.href + "/update", JSON.stringify(data));
}

function GetTestPlan(projectId) {
    let testCasesList = document.getElementById("test-cases");
    testCasesList.replaceChildren();

    let response = post_request(window.location.href + "/get");
    let data = JSON.parse(response)

    let description = document.getElementById("description-text")
    if (data.description != null) {
        description.innerText = data.description
    }

    let records = data.cases;
    for (recordIndex in records) {
        let element = document.createElement("div");
        element.className = "list-item"
        element.draggable = true
        let id = document.createElement("span");
        let name = document.createElement("span");
        let spacer = document.createElement("div")
        spacer.style.flexGrow = 1
        let deleteButton = document.createElement("button")
        deleteButton.innerText = "Удалить"
        deleteButton.className = "hidden"
        deleteButton.onclick = (event) => {
            event.stopPropagation()
            testCasesList.removeChild(element)
            UpdateTestPlatTestCasesList(testCasesList)
        }

        id.innerText = records[recordIndex].id;
        name.innerText = records[recordIndex].name;

        element.appendChild(id);
        element.appendChild(name);
        element.appendChild(spacer);
        element.appendChild(deleteButton);

        let testCaseId = id.innerText;
        element.onclick = () => OpenPage("/project/" + projectId + "/case/" + testCaseId);

        element.addEventListener("dragstart", () => {
            setTimeout(() => { element.classList.add("dragging") }, 0);
        });

        element.addEventListener("dragend", () => {
            element.classList.remove("dragging")

            let data = { "cases": [] }
            for (i = 0; i < testCasesList.children.length; i++) {
                let testCase = testCasesList.children.item(i)
                let id = testCase.children.item(0).innerText
                data["cases"].push({ "id": id })
            }
            let response = post_request(window.location.href + "/update", JSON.stringify(data));
        });

        testCasesList.appendChild(element);
    }

    const initSortableList = (event) => {
        event.preventDefault();
        const draggingItem = document.querySelector(".dragging");
        let siblings = [...testCasesList.querySelectorAll(".list-item:not(.dragging)")];
        let nextSibling = siblings.find(sibling => {
            return event.clientY + window.scrollY < sibling.offsetTop + sibling.offsetHeight / 2;
        });
        testCasesList.insertBefore(draggingItem, nextSibling);
    }
    testCasesList.addEventListener("dragover", initSortableList);
    testCasesList.addEventListener("dragenter", event => event.preventDefault());

    let appendCaseSelect = document.getElementById("append-case-select")

    let testCasesResponse = post_request(
        window.request_url + "/project/" + projectId + "/cases/get",
        JSON.stringify({ "search": "" })
    );

    let testCaseRecords = JSON.parse(testCasesResponse).records;
    for (recordIndex in testCaseRecords) {
        let record = testCaseRecords[recordIndex]
        let appendCaseOption = document.createElement("option")
        appendCaseOption.innerText = record.fields.Name;
        appendCaseOption.__custom__ = {}
        appendCaseOption.__custom__.id = record.fields.Id;
        appendCaseSelect.appendChild(appendCaseOption)
    }

    UpdateTestPlanTags();
}

function Edit(editButton, textElementId, textAreaElementId, fieldName) {
    let textElement = document.getElementById(textElementId)
    let textPlaceholderElement = document.getElementById(textElementId + "-placeholder")
    let textAreaElement = document.getElementById(textAreaElementId)

    if (editButton.editMode == true) {
        editButton.editMode = false
        editButton.style.width = "30px"
        editButton.firstElementChild.style.display = "block"
        editButton.lastElementChild.style.display = "none"

        textElement.innerText = textAreaElement.value

        textElement.style.display = "block"
        textAreaElement.style.display = "none"
        if (textPlaceholderElement != null && textElement.innerText.length == 0) {
            textPlaceholderElement.style.display = "flex"
        }

        let data = {}
        data[fieldName] = textElement.innerText
        post_request(
            window.location.href + "/update",
            JSON.stringify(data)
        )

        return
    }

    editButton.editMode = true
    editButton.style.width = "80px"

    editButton.firstElementChild.style.display = "none"
    editButton.lastElementChild.style.display = "block"

    textAreaElement.value = textElement.innerText

    textElement.style.display = "none"
    textAreaElement.style.display = "block"
    if (textPlaceholderElement != null) {
        textPlaceholderElement.style.display = "none"
    }
}

function AppendTestCase(projectId) {
    let appendCaseSelect = document.getElementById("append-case-select")
    let id = appendCaseSelect.options[appendCaseSelect.selectedIndex].__custom__.id
    let data = id
    let response = post_request(window.location.href + "/case/append", id)
    window.close_dialog('append-dialog', 'append-dialog-overlay')
}
