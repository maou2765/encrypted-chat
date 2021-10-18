import { Subject } from "rxjs";
import { debounceTime } from "rxjs/operators";
import config from "./utils/config";
interface Friend {
  ID: string;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string;
  given_name: string;
  surname: string;
  icon_url: string;
  bio: string;
  email: string;
}
const friendRow = (props: Friend) => `
<button class="row align-items-start friend-row" id="friend-row-${props.ID}">
    <div class="col-sm-3"><img src="https://2sxc.org/Portals/0/adam/Content/4IqBjx3pXEC7a7-fVX2GBQ/Image/github-logo.png?maxwidth=1400&maxheight=990&quality=75" height="72px"/></div>
    <div class="col-sm-9">
        <div class="row align-items-start">
            ${props.given_name}${" " + props.surname}
            </ul>
        </div>
        <div class="row align-items-start">
         ${props.email}
        </div>
        <div class="row align-items-start">
         ${props.bio}
        </div>
    </div>
</button>`;
const setFriendsOption = (friends: Friend[]) => {
  $(".friend-row").remove();
  friends.forEach((friend) => {
    $("#friend-list").append(friendRow(friend));
    const buttonId = `button#friend-row-${friend.ID}`;
    $(buttonId).data("friend", friend);
    $(buttonId).click(
      (
        e: JQuery.ClickEvent<
          HTMLDivElement,
          null,
          HTMLDivElement,
          HTMLDivElement
        >
      ) => {
        console.log("on friend click", $(e.target).data("friend"));
        if ($(buttonId).parent().attr("id") == "added-friends") {
          if ($("#added-friends").children().length == 1) {
            $("#add-friend-subtitle").addClass("hidden");
          }
          $(buttonId).detach().appendTo("#friend-list");
        } else {
          $("#add-friend-subtitle").removeClass("hidden");
          $(buttonId).detach().appendTo("#added-friends");
        }
      }
    );
  });
};
const search = (searchKeyword: string) => {
  const url = new URL(`${config.baseURL}/friends`);
  const params = new URLSearchParams();
  params.append("search", searchKeyword);
  url.search = params.toString();
  fetch(url.toString(), {
    credentials: "same-origin",
  })
    .then((resp) => resp.json())
    .then((resp) => {
      console.log(resp);
      setFriendsOption(resp.friends);
    })
    .catch((e) => console.error(e));
};
$(document).ready(() => {
  const keywordSubject = new Subject<string>();
  const keywordDebounce = keywordSubject.pipe(debounceTime(100));
  keywordDebounce.subscribe((searchKeyword) => {
    search(searchKeyword);
  });
  $("#search-keyword").keyup(
    (
      e: JQuery.KeyUpEvent<
        HTMLInputElement,
        null,
        HTMLInputElement,
        HTMLInputElement
      >
    ) => {
      keywordSubject.next(e.target.value);
    }
  );
  $("#search-keyword").keypress(
    (
      e: JQuery.KeyPressEvent<
        HTMLInputElement,
        null,
        HTMLInputElement,
        HTMLInputElement
      >
    ) => {
      if (e.code == "Enter") {
        search($("#search-keyword").val() as string);
      }
    }
  );
  $("#add-btn").click(
    (
      e: JQuery.ClickEvent<
        HTMLButtonElement,
        null,
        HTMLButtonElement,
        HTMLButtonElement
      >
    ) => {
      const addedFds = [];
      const formData = new FormData();
      $("#added-friends > .friend-row").map(function () {
        const data = $(this).data();
        if (data && data.friend) {
          addedFds.push(data.friend.ID);
          formData.append(`fd[]`, data.friend.ID);
        }
      });
      console.log("addedFds", addedFds);
      if (addedFds.length > 0) {
        const url = new URL(`${config.baseURL}/friends`);
        fetch(url.toString(), {
          method: "POST",
          credentials: "same-origin",
          body: formData,
        })
          .then((resp) => resp.json())
          .then((resp) => {
            console.log(resp);
          })
          .catch((e) => console.error(e));
      }
    }
  );
});
